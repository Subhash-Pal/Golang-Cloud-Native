package main

import (
	"fmt"
	"strings"
)

type PricingStrategy interface {
	Calculate(total float64) float64
}

type RegularPricing struct{}

func (RegularPricing) Calculate(total float64) float64 {
	return total
}

type FestivalDiscount struct{}

func (FestivalDiscount) Calculate(total float64) float64 {
	return total * 0.85
}

type PremiumMemberDiscount struct{}

func (PremiumMemberDiscount) Calculate(total float64) float64 {
	return total * 0.75
}

func NewPricingStrategy(name string) (PricingStrategy, error) {
	switch strings.ToLower(name) {
	case "regular":
		return RegularPricing{}, nil
	case "festival":
		return FestivalDiscount{}, nil
	case "premium":
		return PremiumMemberDiscount{}, nil
	default:
		return nil, fmt.Errorf("unknown pricing strategy: %s", name)
	}
}

type OrderService struct{}

func (OrderService) FinalPrice(strategy PricingStrategy, total float64) float64 {
	return strategy.Calculate(total)
}

func main() {
	service := OrderService{}

	for _, name := range []string{"regular", "festival", "premium", "invalid"} {
		strategy, err := NewPricingStrategy(name)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		finalPrice := service.FinalPrice(strategy, 2000)
		fmt.Printf("%s final price: %.2f\n", name, finalPrice)
	}
}
