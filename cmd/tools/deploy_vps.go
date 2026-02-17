package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	VPS_IP       = "43.134.180.17"
	VPS_USER     = "root"
	VPS_PASSWORD = "u9K-4dF-zrn-2ie"
	REPO_URL     = "https://github.com/imamrisnandar/muslimly-be.git"
	APP_DIR      = "/root/muslimly-be"
)

func main() {
	// 1. Setup SSH Config
	config := &ssh.ClientConfig{
		User: VPS_USER,
		Auth: []ssh.AuthMethod{
			ssh.Password(VPS_PASSWORD),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	fmt.Printf("Connecting to %s...\n", VPS_IP)
	client, err := ssh.Dial("tcp", VPS_IP+":22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer client.Close()
	fmt.Println("Connected! checking VPS info...")

	// 2. Check info
	runCommand(client, "uname -a")
	runCommand(client, "cat /etc/os-release | grep PRETTY_NAME")

	// 3. Install Prerequisites (Smart Detection)
	fmt.Println("\n--- Installing Prerequisites ---")
	installCmd := `
	set -e
	if command -v apt-get &> /dev/null; then
		echo "Detected Debian/Ubuntu"
		apt-get update -qq
		apt-get install -y git curl
	elif command -v yum &> /dev/null; then
		echo "Detected RHEL/CentOS/OpenCloudOS"
		yum install -y git curl
	else
		echo "Unknown Package Manager"
		exit 1
	fi

	# Install Docker
	if ! command -v docker &> /dev/null; then
		echo "Installing Docker..."
		if command -v yum &> /dev/null; then
			yum install -y docker
			systemctl start docker
			systemctl enable docker
			# Compose plugin
			yum install -y docker-compose-plugin || curl -SL https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose
		else
			curl -fsSL https://get.docker.com -o get-docker.sh
			sh get-docker.sh
		fi
	else
		echo "Docker already installed"
	fi
	`
	runCommand(client, installCmd)

	// 4. Setup Repo
	fmt.Println("\n--- Setting up Repository ---")
	// Make sure dir exists or clone
	checkRepoCmd := fmt.Sprintf(`
	if [ -d "%s" ]; then
		echo "Repo exists, pulling..."
		cd %s && git pull || echo "Git pull failed, ignoring..."
	else
		echo "Cloning repo..."
		git clone %s %s
	fi
	mkdir -p %s
	`, APP_DIR, APP_DIR, REPO_URL, APP_DIR, APP_DIR)
	runCommand(client, checkRepoCmd)

	// 5. Upload Secrets
	fmt.Println("\n--- Uploading Config Files ---")
	uploadFile(client, "config.yaml", APP_DIR+"/config.yaml")
	uploadFile(client, "firebase-service-account.json", APP_DIR+"/firebase-service-account.json")

	// 6. Create .env for Docker Compose
	fmt.Println("\n--- Creating .env ---")
	envContent := "DB_USER=postgres\nDB_PASSWORD=muslimly_secure_pass\nDB_NAME=muslimly\n"
	createEnvCmd := fmt.Sprintf("echo '%s' > %s/.env", envContent, APP_DIR)
	runCommand(client, createEnvCmd)

	// 7. Deploy
	fmt.Println("\n--- Deploying Container ---")
	// Check if 'docker compose' (plugin) or 'docker-compose' (standalone) exists
	deployCmd := fmt.Sprintf(`
	cd %s
	if docker compose version >/dev/null 2>&1; then
		docker compose up -d --build
	elif command -v docker-compose >/dev/null 2>&1; then
		docker-compose up -d --build
	else
		echo "Docker Compose not found!"
		exit 1
	fi
	`, APP_DIR)
	runCommand(client, deployCmd)

	fmt.Println("\n✅ Deployment Completed! Check API at http://" + VPS_IP + ":8080")
}

func runCommand(client *ssh.Client, cmd string) {
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	logLen := 50
	if len(cmd) < 50 {
		logLen = len(cmd)
	}
	// Sanitize newlines for log
	logCmd := cmd[:logLen]

	fmt.Printf("Running: %q...\n", logCmd) // Log header
	if err := session.Run(cmd); err != nil {
		log.Printf("Command finished with error: %v", err)
		// Don't exit, try to continue (idempotent)
	}
}

func uploadFile(client *ssh.Client, localPath, remotePath string) {
	// Read Local
	f, err := os.Open(localPath)
	if err != nil {
		log.Fatalf("Failed to open local file %s: %v", localPath, err)
	}
	defer f.Close()

	stat, _ := f.Stat()
	size := stat.Size()

	// New Session for SCP/Cat
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session for upload: %v", err)
	}
	defer session.Close()

	// Use 'cat' on remote to write file
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		io.Copy(w, f)
	}()

	fmt.Printf("Uploading %s (%d bytes)...\n", localPath, size)
	if err := session.Run("cat > " + remotePath); err != nil {
		log.Fatalf("Failed to upload: %v", err)
	}
}
