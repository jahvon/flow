package executable

import (
	"fmt"
	"time"
)

func WithTimeout(timeoutStr string, fn func() error) error {
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return fmt.Errorf("unable to parse timeout duration - %v", err)
	}
	if timeout == 0 {
		return fn()
	}
	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("timeout after %v", timeout)
	}
}
