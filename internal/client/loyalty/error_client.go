package loyalty

import (
	"fmt"
)

type ErrOrderNotRegistered struct {
	OrderNumber string
}

func (e *ErrOrderNotRegistered) Error() string {
	return fmt.Sprintf("order %s not registered in loyalty system", e.OrderNumber)
}

type ErrRateLimitExceeded struct {
	RetryAfter string // Header Retry-After
}

func (e *ErrRateLimitExceeded) Error() string {
	return "rate limit exceeded for loyalty service"
}

type ErrServerError struct {
	StatusCode int
}

func (e *ErrServerError) Error() string {
	return fmt.Sprintf("loyalty service returned server error: %d", e.StatusCode)
}

type ErrUnexpectedStatus struct {
	StatusCode int
	Status     string
}

func (e *ErrUnexpectedStatus) Error() string {
	return fmt.Sprintf("unexpected status from loyalty service: %d %s",
		e.StatusCode, e.Status)
}
