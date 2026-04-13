package sync

import (
	"fmt"
	"math/rand"
	"time"
)

func WithRetry(attempts int, baseDelay time.Duration, fn func() error) error {
	var err error
	delay := baseDelay
	for i := 1; i <= attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if i == attempts {
			break
		}
		jitter := time.Duration(rand.Intn(1001)) * time.Millisecond
		time.Sleep(delay + jitter)
		delay *= 5
	}
	return fmt.Errorf("final failure after %d attempts: %w", attempts, err)
}
