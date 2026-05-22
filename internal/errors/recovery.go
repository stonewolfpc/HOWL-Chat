package errors

import (
	"fmt"
	"runtime/debug"
)

// ---------------------------------------------------------------------------
// Recovery - Panic recovery wrappers for goroutines and critical sections
// Every goroutine that calls into application code MUST use one of these.
// ---------------------------------------------------------------------------

// RecoverFromPanic recovers from a panic, creates a TrackedError, and calls
// the provided handler. Use this at the top of every goroutine.
func RecoverFromPanic(operation string, handler func(*TrackedError)) {
	if r := recover(); r != nil {
		err := New(
			SYSPanicRecovered,
			Critical,
			operation,
			getPanicLocation(),
			fmt.Errorf("panic: %v", r),
		).WithContext("stack", string(debug.Stack()))

		if handler != nil {
			handler(err)
		}
	}
}

// RecoverSilent recovers from a panic without surfacing to user.
// Use for background cleanup goroutines.
func RecoverSilent(operation string, logger *Logger) {
	if r := recover(); r != nil {
		err := New(
			SYSPanicRecovered,
			Medium,
			operation,
			getPanicLocation(),
			fmt.Errorf("panic: %v", r),
		).WithContext("stack", string(debug.Stack()))

		if logger != nil {
			logger.Log(err)
		}
	}
}

// Guard executes fn and recovers from any panic, returning the error.
// Use for critical synchronous operations that must not crash.
func Guard(operation string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = New(
				SYSPanicRecovered,
				Critical,
				operation,
				getPanicLocation(),
				fmt.Errorf("panic: %v", r),
			).WithContext("stack", string(debug.Stack()))
		}
	}()
	return fn()
}

// GuardValue executes fn and recovers from any panic, returning value and error.
func GuardValue[T any](operation string, fn func() (T, error)) (val T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = New(
				SYSPanicRecovered,
				Critical,
				operation,
				getPanicLocation(),
				fmt.Errorf("panic: %v", r),
			).WithContext("stack", string(debug.Stack()))
		}
	}()
	return fn()
}

// WrapError converts a standard error into a TrackedError.
// Use when crossing subsystem boundaries to attach metadata.
func WrapError(code Code, severity Severity, operation, location string, err error) *TrackedError {
	if err == nil {
		return nil
	}
	if te, ok := err.(*TrackedError); ok {
		return te
	}
	return New(code, severity, operation, location, err)
}

// getPanicLocation extracts the file:line where the panic originated.
func getPanicLocation() string {
	return string(debug.Stack())
}
