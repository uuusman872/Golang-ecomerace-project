package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/uuusman872/learning-go-shop/internal/config"
	"github.com/uuusman872/learning-go-shop/internal/models"
	"github.com/uuusman872/learning-go-shop/internal/notifications"
	"github.com/uuusman872/learning-go-shop/internal/providers"
)

func main() {
	log.Println("Starting notification service...")
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	emailConfig := &notifications.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}
	emailNotifier := notifications.NewEmailNotifier(emailConfig)
	awsConfig, err := providers.CreateAwsConfig(ctx, cfg.AWS.S3Endpoint, cfg.AWS.Region)
	if err != nil {
		log.Fatalf("Failed to create aws config: %v", err)
	}
	logger := watermill.NewStdLogger(false, false)
	subscriber, err := sqs.NewSubscriber(sqs.SubscriberConfig{
		AWSConfig: awsConfig,
	}, logger)
	if err != nil {
		log.Fatalf("Failed to create subscriber: %v", err)
	}
	message, err := subscriber.Subscribe(ctx, cfg.AWS.EventQueueName)
	if err != nil {
		subscriber.Close()
		log.Fatalf("Failed to subscribe: %v", err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Notification service started")
	for {
		select {
		case msg := <-message:
			if err := processMessage(msg, emailNotifier); err != nil {
				log.Printf("Failed to process message: %v", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		case <-sigChan:
			log.Println("Shutting down notification service...")
			subscriber.Close()
			return
		}
	}
}

func processMessage(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	eventType := msg.Metadata.Get("event_type")
	switch eventType {
	case notifications.UserLoggedIn:
		return handleUserLoggedIn(msg, emailNotifier)
	default:
		log.Printf("UnKnown event type: %s\n", eventType)
		return nil
	}
}

func handleUserLoggedIn(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	var user models.User
	if err := json.Unmarshal(msg.Payload, &user); err != nil {
		return err
	}
	userName := user.FirstName + " " + user.LastName
	if userName != " " {
		userName = "User"
	}
	log.Printf("Sending login notification to %s", user.Email)
	return emailNotifier.SendLoginNotification(user.Email, userName)
}
