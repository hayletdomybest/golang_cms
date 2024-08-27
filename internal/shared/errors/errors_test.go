package errors

import (
	"errors"
	"testing"
)

func TestErrorRegistration(t *testing.T) {
	tests := []struct {
		err       *Error
		codespace string
		code      uint32
		desc      string
	}{
		{ErrUnexpected, "ERR_UNEXPECTED", 1, "An unexpected error occurred"},
		{ErrInternal, "ERR_INTERNAL", 2, "Internal server error"},
		{ErrNotFound, "ERR_NOT_FOUND", 3, "Resource not found"},
		{ErrInvalidParams, "ERR_INVALID_Params", 4, "Invalid input provided"},
		{ErrUnauthorized, "ERR_UNAUTHORIZED", 5, "Unauthorized access"},
		{ErrAlreadyExist, "ERR_ALREADY_EXISTS", 6, "Resource already exists"},
	}

	for _, tt := range tests {
		if tt.err.codespace != tt.codespace || tt.err.code != tt.code || tt.err.desc != tt.desc {
			t.Errorf("Expected %v, got %v", tt, tt.err)
		}
	}
}

func TestErrorWrapping(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := ErrInternal.Wrap(originalErr)

	if wrappedErr.cause != originalErr {
		t.Errorf("Expected cause to be original error, got %v", wrappedErr.cause)
	}

	if wrappedErr.codespace != ErrInternal.codespace || wrappedErr.code != ErrInternal.code || wrappedErr.desc != ErrInternal.desc {
		t.Errorf("Wrapped error does not match original error properties")
	}
}

func TestErrorWithMessage(t *testing.T) {
	customMessage := "custom message"
	errWithMsg := ErrNotFound.WithMessage(customMessage)

	if errWithMsg.desc != customMessage {
		t.Errorf("Expected custom message '%s', got '%s'", customMessage, errWithMsg.desc)
	}
}

func TestIsFunction(t *testing.T) {
	wrappedErr := ErrUnauthorized.Wrap(ErrInternal)

	if !ErrUnauthorized.Is(wrappedErr) {
		t.Errorf("Expected wrapped error to match ErrUnauthorized")
	}

	if ErrInternal.Is(wrappedErr) {
		t.Errorf("Did not expect wrapped error to match ErrInternal")
	}

	customMessage := "custom message"
	errWithMsg := ErrNotFound.WithMessage(customMessage)

	if !errWithMsg.Is(ErrNotFound) {
		t.Errorf("Expected error to match ErrNotFound")
	}
}

func TestUniqueErrorRegistration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on duplicate error registration, but it did not occur")
		}
	}()

	// This should cause a panic since ErrUnexpected with code 1 is already registered
	Register("ERR_UNEXPECTED", 1, "This should panic")
}
