package service

import (
	"context"
	"fmt"
	"log"
	"muslimly-be/internal/features/notification/dto"
	"muslimly-be/internal/features/notification/model"
	"muslimly-be/internal/features/notification/repository"
	"muslimly-be/pkg/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type notificationService struct {
	repo   repository.DeviceRepository
	config *config.Config
	fcm    *messaging.Client
}

func NewNotificationService(repo repository.DeviceRepository, config *config.Config) NotificationService {
	// Initialize Firebase
	opt := option.WithCredentialsFile(config.Notification.FirebaseCredentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("WARNING: Failed to init Firebase App: %v", err)
		return &notificationService{repo, config, nil}
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Printf("WARNING: Failed to init FCM Client: %v", err)
		return &notificationService{repo, config, nil}
	}

	log.Println("Firebase FCM Client Initialized")
	return &notificationService{repo, config, client}
}

func (s *notificationService) RegisterDevice(userID string, req dto.RegisterDeviceRequest) error {
	var uid *uuid.UUID

	// Only parse if userID is not empty (Logged In)
	if userID != "" {
		parsed, err := uuid.Parse(userID)
		if err != nil {
			return err
		}
		uid = &parsed
	}

	device := &model.UserDevice{
		UserID:   uid,
		FCMToken: req.FCMToken,
		Platform: req.Platform,
		// Device Metadata
		DeviceModel:     req.DeviceModel,
		DeviceOSVersion: req.DeviceOSVersion,
		AppVersion:      req.AppVersion,
		CountryCode:     req.CountryCode,
		Timezone:        req.Timezone,
	}

	return s.repo.UpsertDevice(device)
}

func (s *notificationService) SendDailyReminder() error {
	if s.fcm == nil {
		log.Println("FCM Client not initialized, skipping reminder")
		return fmt.Errorf("FCM client not initialized")
	}

	devices, err := s.repo.GetAllDevices()
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return nil
	}

	var guestTokens []string
	var userTokens []string
	var userIDs []string

	for _, d := range devices {
		if d.UserID == nil {
			guestTokens = append(guestTokens, d.FCMToken)
		} else {
			userTokens = append(userTokens, d.FCMToken)
			userIDs = append(userIDs, d.UserID.String())
		}
	}

	// 1. Send to GUESTS (Blind Blast)
	if len(guestTokens) > 0 {
		log.Printf("DEBUG: Sending to %d GUESTS", len(guestTokens))
		successCount := 0
		failureCount := 0

		for _, token := range guestTokens {
			msg := &messaging.Message{
				Token: token,
				Notification: &messaging.Notification{
					Title: "Waktunya Mengaji 📖",
					Body:  "Assalamualaikum, mari sempatkan membaca Al-Quran hari ini.",
				},
				Data: map[string]string{"type": "daily_reminder"},
			}

			_, err := s.fcm.Send(context.Background(), msg)
			if err != nil {
				failureCount++
				log.Printf("Guest Send Error for token %s...: %v", token[:10], err)
			} else {
				successCount++
			}
		}
		log.Printf("Guest Reminder: %d success, %d failure", successCount, failureCount)
	} else {
		log.Println("DEBUG: No Guest Tokens found")
	}

	// 2. Send to USERS
	if len(userTokens) > 0 {
		log.Printf("DEBUG: Sending to %d USERS", len(userTokens))
		successCount := 0
		failureCount := 0

		for _, token := range userTokens {
			msg := &messaging.Message{
				Token: token,
				Notification: &messaging.Notification{
					Title: "Sudahkah Mengaji? 🌙",
					Body:  "Jangan lupa target harianmu. Sempatkan membaca Al-Quran ya.",
				},
				Data: map[string]string{"type": "daily_reminder"},
			}

			_, err := s.fcm.Send(context.Background(), msg)
			if err != nil {
				failureCount++
				log.Printf("User Send Error for token %s...: %v", token[:10], err)
			} else {
				successCount++
			}
		}
		log.Printf("User Reminder: %d success, %d failure", successCount, failureCount)
	} else {
		log.Println("DEBUG: No User Tokens found")
	}

	return nil
}
