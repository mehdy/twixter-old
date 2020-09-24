package interactors

import (
	"fmt"
)

type interactorError struct {
	cause  error
	detail string
}

func newError(cause error, detail string) *interactorError {
	return &interactorError{cause: cause, detail: detail}
}

func (e *interactorError) Error() string {
	return fmt.Sprintf("%s: %s", e.cause, e.detail)
}

func (e *interactorError) Unwrap() error {
	return e.cause
}
