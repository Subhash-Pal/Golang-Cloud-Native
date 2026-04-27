package main

import "fmt"

// Step 1: Define an Interface
// PaymentMethod interface defines a common method for all payment types
type PaymentMethod interface {
	Pay(amount float64) string
}

// Step 2: Create Concrete Implementations
// CreditCard struct implementing PaymentMethod
type CreditCard struct{}

func (c *CreditCard) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using Credit Card", amount)
}

// PayPal struct implementing PaymentMethod
type PayPal struct{}

func (p *PayPal) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using PayPal", amount)
}

// Step 3: Implement the Factory Method
// GetPaymentMethod is the Factory Method that returns an appropriate payment method
func GetPaymentMethod(method string) (PaymentMethod, error) {
	switch method {
	case "credit":
		return &CreditCard{}, nil
	case "paypal":
		return &PayPal{}, nil
	default:
		return nil, fmt.Errorf("payment method %s not recognized", method)
	}
}

// Step 4: Use the Factory Method
func main() {
	paymentMethod, err := GetPaymentMethod("credit")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(paymentMethod.Pay(100.50))

	paymentMethod, err = GetPaymentMethod("paypal")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(paymentMethod.Pay(200.75))
}
