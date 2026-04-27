package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

var ErrPermanentFailure = errors.New("permanent failure")

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	service := &PaymentService{
		failUntilAttempt: 3,
		rng:              rng,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := Retry(ctx, RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    4 * time.Second,
		Jitter:      250 * time.Millisecond,
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
	Jitter      time.Duration
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

		if errors.Is(err, ErrPermanentFailure) {
			return fmt.Errorf("stop retrying because the error is permanent: %w", err)
		}

		lastErr = err
		log.Printf("Attempt %d failed: %v", attempt, err)

		if attempt == cfg.MaxAttempts {
			break
		}

		delay := nextDelay(cfg, attempt)
		log.Printf("Waiting %s before retrying", delay)

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

func nextDelay(cfg RetryConfig, attempt int) time.Duration {
	delay := float64(cfg.BaseDelay) * math.Pow(2, float64(attempt-1))
	if time.Duration(delay) > cfg.MaxDelay {
		delay = float64(cfg.MaxDelay)
	}

	if cfg.Jitter > 0 {
		jitter := rand.Int63n(int64(cfg.Jitter))
		delay += float64(time.Duration(jitter))
	}

	return time.Duration(delay)
}

type PaymentService struct {
	currentAttempt   int
	failUntilAttempt int
	rng              *rand.Rand
}

func (p *PaymentService) ProcessPayment() error {
	p.currentAttempt++

	if p.currentAttempt < p.failUntilAttempt {
		return errors.New("temporary network error")
	}

	if p.rng.Intn(10) == 0 {
		return ErrPermanentFailure
	}

	return nil
}
