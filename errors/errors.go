package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

//Reference article:
//https://hackernoon.com/golang-handling-errors-gracefully-8e27f1db729f

type customError struct {
	errType       ErrorType
	originalError error
	context       errorContext
	stacktrace    bool
}

func (err customError) Error() string {
	if err.stacktrace {
		return err.Stacktrace()
	}
	return err.originalError.Error()
}

func (err customError) Stacktrace() string {
	return fmt.Sprintf("%+v\n", err.originalError)
}

// New creates a no type error
func New(msg string, stacktrace bool) error {
	return customError{errType: Unknown, originalError: errors.New(msg), stacktrace: stacktrace}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return customError{errType: Unknown, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap wrans an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf wraps an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errType:       customErr.errType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}

	return customError{errType: Unknown, originalError: wrappedError}
}

// Get Stacktrace of error
func Stack(err error) string {
	if customErr, ok := err.(customError); ok {
		return fmt.Sprintf("%+v\n", customErr.originalError)
	}
	return fmt.Sprintf("%+v\n", errors.WithStack(err))
}
