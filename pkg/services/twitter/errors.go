package twitter

import "fmt"

type apiError struct {
	cause  error
	detail string
}

func newError(cause error, detail string) *apiError {
	return &apiError{cause: cause, detail: detail}
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s: %s", e.cause, e.detail)
}

func (e *apiError) Unwrap() error {
	return e.cause
}
