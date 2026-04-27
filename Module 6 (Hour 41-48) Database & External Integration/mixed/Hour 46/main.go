package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	service := &PaymentService{failUntilAttempt: 3}

	err := Retry(ctx, RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   1 * time.Second,
		MaxDelay:    5 * time.Second,
	}, func(attempt int) error {
		log.Printf("Attempt %d to process payment", attempt)
		return service.ProcessPayment()
	})
	if err != nil {
		log.Fatalf("operation failed after retries: %v", err)
	}

	log.Println("Payment processed successfully after retry logic")
}

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func Retry(ctx context.Context, cfg RetryConfig, operation func(attempt int) error) error {
	if cfg.MaxAttempts < 1 {
		return fmt.Errorf("max attempts must be at least 1")
	}
	if cfg.BaseDelay <= 0 {
		return fmt.Errorf("base delay must be greater than 0")
	}
	if cfg.MaxDelay < cfg.BaseDelay {
		return fmt.Errorf("max delay must be greater than or equal to base delay")
	}

	var lastErr error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err := operation(attempt)
		if err == nil {
			return nil
		}

		lastErr = err
		log.Printf("Attempt %d failed: %v", attempt, err)

		if attempt == cfg.MaxAttempts {
			break
		}

		delay := nextDelay(cfg.BaseDelay, cfg.MaxDelay, attempt)
		log.Printf("Waiting %s before retrying...", delay)

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("retry canceled: %w", ctx.Err())
		case <-timer.C:
		}
	}

	return fmt.Errorf("all retry attempts failed: %w", lastErr)
}

func nextDelay(baseDelay, maxDelay time.Duration, attempt int) time.Duration {
	delay := float64(baseDelay) * math.Pow(2, float64(attempt-1))
	if time.Duration(delay) > maxDelay {
		return maxDelay
	}
	return time.Duration(delay)
}

type PaymentService struct {
	currentAttempt   int
	failUntilAttempt int
}

func (p *PaymentService) ProcessPayment() error {
	p.currentAttempt++

	if p.currentAttempt < p.failUntilAttempt {
		return errors.New("temporary network error")
	}

	return nil
}
