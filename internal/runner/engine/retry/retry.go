package retry

import (
	"fmt"
	"time"
)

type Stats struct {
	Attempts int
	Failures int
}

type Handler struct {
	maxRetries  int
	backoffTime time.Duration
	stats       Stats
}

func NewRetryHandler(maxRetries int, backoffTime time.Duration) *Handler {
	return &Handler{
		maxRetries:  maxRetries,
		backoffTime: backoffTime,
		stats:       Stats{},
	}
}

func (h *Handler) Execute(operation func() error) error {
	var lastErr error

	for h.stats.Attempts <= h.maxRetries {
		h.stats.Attempts++

		if err := operation(); err != nil {
			h.stats.Failures++
			lastErr = err

			if !h.Retryable() {
				break
			}

			if h.backoffTime > 0 {
				time.Sleep(h.backoffTime)
			}

			continue
		}

		return nil
	}

	if h.maxRetries <= 0 {
		return lastErr
	}
	return fmt.Errorf("execution failed after %d attempts. Last error: %w", h.stats.Attempts, lastErr)
}

func (h *Handler) GetStats() Stats {
	return h.stats
}

func (h *Handler) Retryable() bool {
	return h.stats.Attempts <= h.maxRetries
}

func (h *Handler) Reset() {
	h.stats = Stats{}
}
