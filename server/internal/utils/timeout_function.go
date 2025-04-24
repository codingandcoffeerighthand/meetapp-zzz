package utils

import (
	"context"
	"errors"
)

func WithContext(ctx context.Context, f func() error) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- f()
	}()

	select {
	case <-ctx.Done():
		return errors.New("timeout error")
	case err := <-errChan:
		return err
	}
}
