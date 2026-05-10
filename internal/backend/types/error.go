/**
 * Error Types
 * 
 * This package defines all error-related types for the chat system.
 * Provides structured error handling with specific error codes and messages.
 * 
 * @package types
 */

package types

import "fmt"

// ErrorCode represents a specific error type
type ErrorCode string

const (
	// ErrorCodeModelLoad represents a model loading error
	ErrorCodeModelLoad ErrorCode = "MODEL_LOAD_ERROR"
	// ErrorCodeModelUnload represents a model unloading error
	ErrorCodeModelUnload ErrorCode = "MODEL_UNLOAD_ERROR"
	// ErrorCodeModelNotFound represents a model not found error
	ErrorCodeModelNotFound ErrorCode = "MODEL_NOT_FOUND"
	// ErrorCodeModelIncompatible represents an incompatible model error
	ErrorCodeModelIncompatible ErrorCode = "MODEL_INCOMPATIBLE"
	// ErrorCodeInference represents an inference error
	ErrorCodeInference ErrorCode = "INFERENCE_ERROR"
	// ErrorCodeContext represents a context error
	ErrorCodeContext ErrorCode = "CONTEXT_ERROR"
	// ErrorCodeInvalidInput represents an invalid input error
	ErrorCodeInvalidInput ErrorCode = "INVALID_INPUT"
	// ErrorCodeMemory represents a memory allocation error
	ErrorCodeMemory ErrorCode = "MEMORY_ERROR"
	// ErrorCodeIO represents an I/O error
	ErrorCodeIO ErrorCode = "IO_ERROR"
	// ErrorCodeInternal represents an internal system error
	ErrorCodeInternal ErrorCode = "INTERNAL_ERROR"
)

// Error represents a structured error with code and message
type Error struct {
	Code    ErrorCode `json:"code"`    // Error code
	Message string    `json:"message"` // Human-readable error message
	Details string    `json:"details,omitempty"` // Additional error details
	Cause   error     `json:"-"`       // Underlying error (not serialized
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error
func (e *Error) Unwrap() error {
	return e.Cause
}

// NewError creates a new error with the given code and message
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithDetails creates a new error with code, message, and details
func NewErrorWithDetails(code ErrorCode, message, details string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// WrapError wraps an existing error with additional context
func WrapError(code ErrorCode, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsModelLoadError checks if an error is a model load error
func IsModelLoadError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrorCodeModelLoad
	}
	return false
}

// IsModelUnloadError checks if an error is a model unload error
func IsModelUnloadError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrorCodeModelUnload
	}
	return false
}

// IsModelNotFoundError checks if an error is a model not found error
func IsModelNotFoundError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrorCodeModelNotFound
	}
	return false
}

// IsInferenceError checks if an error is an inference error
func IsInferenceError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrorCodeInference
	}
	return false
}

// IsMemoryError checks if an error is a memory error
func IsMemoryError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrorCodeMemory
	}
	return false
}

// Common error constructors for convenience
func ErrModelLoad(message string) *Error {
	return NewError(ErrorCodeModelLoad, message)
}

func ErrModelUnload(message string) *Error {
	return NewError(ErrorCodeModelUnload, message)
}

func ErrModelNotFound(message string) *Error {
	return NewError(ErrorCodeModelNotFound, message)
}

func ErrModelIncompatible(message string) *Error {
	return NewError(ErrorCodeModelIncompatible, message)
}

func ErrInference(message string) *Error {
	return NewError(ErrorCodeInference, message)
}

func ErrContext(message string) *Error {
	return NewError(ErrorCodeContext, message)
}

func ErrInvalidInput(message string) *Error {
	return NewError(ErrorCodeInvalidInput, message)
}

func ErrMemory(message string) *Error {
	return NewError(ErrorCodeMemory, message)
}

func ErrIO(message string) *Error {
	return NewError(ErrorCodeIO, message)
}

func ErrInternal(message string) *Error {
	return NewError(ErrorCodeInternal, message)
}
