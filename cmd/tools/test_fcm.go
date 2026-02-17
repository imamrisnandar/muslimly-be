package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"muslimly-be/pkg/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type DeviceTest struct {
	FCMToken string `json:"fcm_token"`
}

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()
	fmt.Printf("Config Loaded: %s\n", cfg.Notification.FirebaseCredentialsFile)

	// 2. Load Token
	file, err := os.ReadFile("device_test.json")
	if err != nil {
		log.Fatalf("Error reading device_test.json: %v", err)
	}
	var device DeviceTest
	if err := json.Unmarshal(file, &device); err != nil {
		log.Fatalf("Error parsing json: %v", err)
	}
	fmt.Printf("Testing Token: %s...\n", device.FCMToken[:20])

	// 3. Init Firebase
	opt := option.WithCredentialsFile(cfg.Notification.FirebaseCredentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Firebase Init Failed: %v", err)
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("FCM Client Factory Failed: %v", err)
	}

	// 4. Send Message
	msg := &messaging.Message{
		Token: device.FCMToken,
		Notification: &messaging.Notification{
			Title: "Test Direct FCM",
			Body:  "This is a direct test from cmd/tools/test_fcm.go",
		},
	}

	fmt.Println("Sending message...")
	response, err := client.Send(context.Background(), msg)
	if err != nil {
		log.Fatalf("FCM Send FAILED: %v", err)
	}

	fmt.Printf("FCM Send SUCCESS! Message ID: %s\n", response)
}
