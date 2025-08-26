package token

import (
	"fmt"
)

type UnexpectedSigningMethodError struct {
	Algorithm interface{}
}

func (e *UnexpectedSigningMethodError) Error() string {
	return fmt.Sprintf("unexpected signing method: %v", e.Algorithm)
}
