package postgres

import "fmt"

type dbError struct {
	cause  error
	detail string
}

func newError(cause error, detail string) *dbError {
	return &dbError{cause: cause, detail: detail}
}

func (e *dbError) Error() string {
	return fmt.Sprintf("%s: %s", e.cause, e.detail)
}

func (e *dbError) Unwrap() error {
	return e.cause
}
