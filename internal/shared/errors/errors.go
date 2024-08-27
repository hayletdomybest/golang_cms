package errors

import (
	"fmt"
	"sync"
)

var (
	ErrPanic         = Register("ERR_PANIC", 1, "panic error")
	ErrUnexpected    = Register("ERR_UNEXPECTED", 2, "An unexpected error occurred")
	ErrInternal      = Register("ERR_INTERNAL", 3, "Internal server error")
	ErrNotFound      = Register("ERR_NOT_FOUND", 4, "Resource not found")
	ErrInvalidParams = Register("ERR_INVALID_Params", 5, "Invalid input provided")
	ErrUnauthorized  = Register("ERR_UNAUTHORIZED", 6, "Unauthorized access")
	ErrAlreadyExist  = Register("ERR_ALREADY_EXISTS", 7, "Resource already exists")
)

// Error represents a custom error type with additional metadata.
type Error struct {
	codespace string
	code      uint32
	desc      string
	cause     error
}

// Error method to satisfy the error interface.
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s:%d] %s: %v", e.codespace, e.code, e.desc, e.cause)
	}
	return fmt.Sprintf("[%s:%d] %s", e.codespace, e.code, e.desc)
}

// Wrap adds context to the error.
func (e *Error) Wrap(err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		codespace: e.codespace,
		code:      e.code,
		desc:      e.desc,
		cause:     err,
	}
}

// Wrap adds context to the error.
func (e *Error) WithMessagef(format string, args ...interface{}) *Error {
	return e.WithMessage(fmt.Sprintf(format, args...))
}

func (e *Error) WithMessage(msg string) *Error {
	return &Error{
		codespace: e.codespace,
		code:      e.code,
		desc:      msg,
		cause:     e.cause,
	}
}

// Is checks if the error matches the target custom error.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.codespace == t.codespace && e.code == t.code
}

// Register error code registry
var (
	errorRegistry = make(map[string]*Error)
	mutex         sync.Mutex
)

// Register creates a new error with a unique code.
func Register(codespace string, code uint32, description string) *Error {
	err := New(codespace, code, description)

	mutex.Lock()
	defer mutex.Unlock()

	key := fmt.Sprintf("%s:%d", codespace, code)
	if _, exists := errorRegistry[key]; exists {
		panic(fmt.Sprintf("error with code %d in codespace %s is already registered", code, codespace))
	}
	errorRegistry[key] = err

	return err
}

// New creates a new error using the Register function.
func New(codespace string, code uint32, desc string) *Error {
	return &Error{codespace: codespace, code: code, desc: desc}
}
