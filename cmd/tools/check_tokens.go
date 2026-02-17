package main

import (
	"fmt"
	"log"

	"muslimly-be/internal/features/auth_user/domain"
	"muslimly-be/internal/features/notification/model"
	"muslimly-be/pkg/config"
	"muslimly-be/pkg/database"
)

func main() {
	// 1. Load Config & DB
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)

	// 2. Fetch Users
	var users []domain.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("DB Error Users: %v", err)
	}

	// 3. Fetch Devices
	var devices []model.UserDevice
	if err := db.Find(&devices).Error; err != nil {
		log.Fatalf("DB Error Devices: %v", err)
	}

	fmt.Printf("\n--- FCM TOKEN REPORT ---\n")
	fmt.Printf("Total Registered Users: %d\n", len(users))
	fmt.Printf("Total Registered Devices: %d\n", len(devices))

	// Map Devices by UserID
	userDeviceMap := make(map[string]bool)
	guestCount := 0

	for _, d := range devices {
		if d.UserID != nil {
			userDeviceMap[d.UserID.String()] = true
		} else {
			guestCount++
		}
	}

	fmt.Printf("Guest Devices (No Login): %d\n", guestCount)
	fmt.Printf("------------------------\n")

	usersWithToken := 0
	for _, u := range users {
		hasToken := userDeviceMap[u.ID.String()]
		status := "❌ NO TOKEN"
		if hasToken {
			status = "✅ HAS TOKEN"
			usersWithToken++
		}
		// Print only first 10 characters of ID/Name for brevity
		fmt.Printf("[%s] %-20s : %s\n", u.ID.String()[:8], u.Name, status)
	}

	fmt.Printf("------------------------\n")
	fmt.Printf("Summary: %d/%d Users have FCM Tokens.\n", usersWithToken, len(users))
}
