package main

import "fmt"

type Sender interface {
	Send(to, message string) error
}

type Logger interface {
	Info(message string)
}

type EmailSender struct{}

func (EmailSender) Send(to, message string) error {
	fmt.Printf("EMAIL to %s: %s\n", to, message)
	return nil
}

type SMSSender struct{}

func (SMSSender) Send(to, message string) error {
	fmt.Printf("SMS to %s: %s\n", to, message)
	return nil
}

type ConsoleLogger struct{}

func (ConsoleLogger) Info(message string) {
	fmt.Println("LOG:", message)
}

type NotificationService struct {
	sender Sender
	logger Logger
}

func NewNotificationService(sender Sender, logger Logger) *NotificationService {
	return &NotificationService{
		sender: sender,
		logger: logger,
	}
}

func (s *NotificationService) Notify(to, message string) error {
	s.logger.Info("sending notification")
	if err := s.sender.Send(to, message); err != nil {
		return err
	}
	s.logger.Info("notification sent")
	return nil
}

func main() {
	logger := ConsoleLogger{}

	emailService := NewNotificationService(EmailSender{}, logger)
	_ = emailService.Notify("user@example.com", "Welcome by email")

	fmt.Println()

	smsService := NewNotificationService(SMSSender{}, logger)
	_ = smsService.Notify("+91-9999999999", "Welcome by SMS")
}
