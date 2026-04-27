package main

import "fmt"

type PaymentProcessor interface {
	Pay(amount float64) error
}

type LegacyGateway struct{}

func (LegacyGateway) MakePayment(amountInPaise int) bool {
	fmt.Printf("legacy gateway processed %d paise\n", amountInPaise)
	return true
}

type LegacyGatewayAdapter struct {
	legacy LegacyGateway
}

func (a LegacyGatewayAdapter) Pay(amount float64) error {
	ok := a.legacy.MakePayment(int(amount * 100))
	if !ok {
		return fmt.Errorf("payment failed")
	}
	return nil
}
///// ...

type Observer interface {
	Update(message string)
}

type EmailNotifier struct{}

func (EmailNotifier) Update(message string) {
	fmt.Println("email notifier:", message)
}

type AuditLogger struct{}

func (AuditLogger) Update(message string) {
	fmt.Println("audit logger:", message)
}

type PaymentService struct {
	processor PaymentProcessor
	observers []Observer
}

func NewPaymentService(processor PaymentProcessor) *PaymentService {
	return &PaymentService{processor: processor}
}

func (s *PaymentService) Attach(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *PaymentService) notifyAll(message string) {
	for _, observer := range s.observers {
		observer.Update(message)
	}
}

func (s *PaymentService) CompletePayment(amount float64) error {
	if err := s.processor.Pay(amount); err != nil {
		return err
	}
	s.notifyAll(fmt.Sprintf("payment of %.2f completed", amount))
	return nil
}

func main() {
	service := NewPaymentService(LegacyGatewayAdapter{legacy: LegacyGateway{}})
	service.Attach(EmailNotifier{})
	service.Attach(AuditLogger{})

	if err := service.CompletePayment(1499.50); err != nil {
		fmt.Println("error:", err)
	}
}

