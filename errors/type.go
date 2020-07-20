package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	Success              ErrorType = 1
	Unknown              ErrorType = 0
	NotFound             ErrorType = -1
	Empty                ErrorType = -2
	InternalServerError  ErrorType = -3
	BadRequest           ErrorType = -4
	DuplicateError       ErrorType = -5
	Unauthorized         ErrorType = -6
	ParseError           ErrorType = -7
	Forbidden            ErrorType = -8
	CacheGetError        ErrorType = -9
	SerializationError   ErrorType = -10
	DeserializationError ErrorType = -11
	CacheSetError        ErrorType = -12
	CacheRemoveError     ErrorType = -13
)

type ErrorType int

func (errType ErrorType) New(msg string) error {
	return customError{
		errType:       errType,
		originalError: errors.New(msg),
	}
}

// New creates a new customError with formatted message
func (errType ErrorType) Newf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)

	return customError{
		errType:       errType,
		originalError: err,
	}
}

// Wrap creates a new wrapped error
func (errType ErrorType) Wrap(err error, msg string) error {
	return errType.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (errType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)

	return customError{
		errType:       errType,
		originalError: newErr,
	}
}
