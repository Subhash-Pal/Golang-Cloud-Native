package main

import "fmt"

// Step 1: Define the PaymentMethod Interface
// PaymentMethod interface defines a common method for all payment types
type PaymentMethod interface {
	Pay(amount float64) string
}

// Step 2: Implement Concrete Payment Methods
// CreditCard payment implementation (Visa & MasterCard)
type CreditCard struct {
	cardType string
}

func (c *CreditCard) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using %s Credit Card", amount, c.cardType)
}

// PayPal payment implementation
type PayPal struct{}

func (p *PayPal) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using PayPal", amount)
}

// Step 3: Define the Abstract Factory Interface
// PaymentFactory defines a factory for creating payments
type PaymentFactory interface {
	CreatePayment() PaymentMethod
}

// Step 4: Implement Concrete Factories
// VisaFactory creates Visa CreditCard
type VisaFactory struct{}

func (v *VisaFactory) CreatePayment() PaymentMethod {
	return &CreditCard{cardType: "Visa"}
}

// MasterCardFactory creates MasterCard CreditCard
type MasterCardFactory struct{}

func (m *MasterCardFactory) CreatePayment() PaymentMethod {
	return &CreditCard{cardType: "MasterCard"}
}

// PayPalFactory creates PayPal payments
type PayPalFactory struct{}

func (p *PayPalFactory) CreatePayment() PaymentMethod {
	return &PayPal{}
}

// Step 5: Implement the Factory Selector
// GetPaymentFactory returns the appropriate payment factory
func GetPaymentFactory(method string) (PaymentFactory, error) {
	switch method {
	case "visa":
		return &VisaFactory{}, nil
	case "mastercard":
		return &MasterCardFactory{}, nil
	case "paypal":
		return &PayPalFactory{}, nil
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", method)
	}
}

// Step 6: Use the Abstract Factory
func main() {
	// Get Visa Payment Factory
	visaFactory, err := GetPaymentFactory("visa")
	if err != nil {
		fmt.Println(err)
		return
	}
	visaPayment := visaFactory.CreatePayment()
	fmt.Println(visaPayment.Pay(100.50))

	// Get MasterCard Payment Factory
	mastercardFactory, err := GetPaymentFactory("mastercard")
	if err != nil {
		fmt.Println(err)
		return
	}
	mastercardPayment := mastercardFactory.CreatePayment()
	fmt.Println(mastercardPayment.Pay(200.75))

	// Get PayPal Payment Factory
	paypalFactory, err := GetPaymentFactory("paypal")
	if err != nil {
		fmt.Println(err)
		return
	}
	paypalPayment := paypalFactory.CreatePayment()
	fmt.Println(paypalPayment.Pay(300.00))
}
