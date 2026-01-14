package errwrap

import "fmt"

// ErrorWrapper with a code, message, and internal error description.
type ErrorWrapper struct {
	code     int
	message  any
	internal error
}

// New constructor ErrorWrapper
func New(options ...Option) *ErrorWrapper {
	e := &ErrorWrapper{
		code:    defaultCode,
		message: defaultMessage,
	}
	// options processing
	for _, o := range options {
		o(e)
	}
	return e
}

// Error get text error
func (e *ErrorWrapper) Error() string {
	if e.code == -1 {
		return fmt.Sprintf("%v", e.message)
	}
	if e.message == nil {
		return fmt.Sprintf("code=%v", e.code)
	}
	return fmt.Sprintf("code=%v, message=%v", e.code, e.message)
}

// Code error status code
func (e *ErrorWrapper) Code() int {
	return e.code
}

// Message short status message
func (e *ErrorWrapper) Message() string {
	return fmt.Sprintf("%v", e.message)
}

// Internal error description
func (e *ErrorWrapper) Internal() error {
	return e.internal
}

// Detailed full error description with code, message and internal
func (e *ErrorWrapper) Detailed() error {
	if e.code == -1 && e.message == nil {
		return fmt.Errorf("internal=%w", e.Internal())
	}
	return fmt.Errorf("%s, internal=%w", e.Error(), e.Internal())
}
