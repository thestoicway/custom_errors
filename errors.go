package customerrors

import (
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
)

type ErrorCode = int

// Non-categorized errors (1-1000)
const (
	// ErrUnknown is an unknown error.
	// This should be used only as a fallback error when something unexpected happens.
	// For example, if the error is not handled in the application, like 500 Internal Server Error.
	ErrUnknown ErrorCode = 1
	// ErrWrongInput is returned when wrong input is provided.
	// For example, if body of the request is not valid JSON or required fields are missing.
	ErrWrongInput ErrorCode = 2
)

// User-related errors (1001-2000)
const (
	// ErrWrongCredentials is returned when the user provides wrong credentials.
	// For example, if the user provides wrong email or password during login
	ErrWrongCredentials ErrorCode = 1001
	// ErrDuplicateEmail is returned when the user tries
	// to register with an email that is already in the system.
	ErrDuplicateEmail ErrorCode = 1002
	// ErrUnauthorized is returned when the user is not authorized to perform the operation.
	// For example, if the user tries to access a resource that is not his.
	ErrUnauthorized ErrorCode = 1003
)

// CustomError represents a custom error.
// It can be used to provide additional information about the error.
// The main objective of this struct is to provide a unified error response
// for all the errors that can occur in the application.
//
// This enables client applications to handle errors in a unified way
// and display a user-friendly error message to the user.
type CustomError struct {
	// Code is a machine-readable error code.
	// It can be used by client applications to handle errors in a unified way.
	// For example, if the error code is ErrWrongCredentials,
	// then the client application can display a message like
	// "Wrong credentials, either email or password is wrong".
	Code int

	// Err is an initial error object that can be used to get the error message
	// and stack trace.
	Err error

	// Details can be used to provide additional information about the error.
	// For example, if user is banned from OTP, then details can contain
	// the date when the ban will be lifted.
	Details interface{}
}

// Error returns the error message.
// It is needed to implement the error interface.
func (e *CustomError) Error() string {
	return fmt.Sprintf("ErrorCode: %d,  %v", e.Code, e.Err.(*errors.Error).ErrorStack())
}

// Unwrap returns the initial error object.
// It is needed to implement the errors.Wrapper interface.
func (e *CustomError) Unwrap() error {
	return e.Err
}

// / Is returns true if the target error is of the same type as the receiver.
func (e *CustomError) Is(target error) bool {
	t, ok := target.(*CustomError)
	if !ok {
		return false
	}

	return e.Code == t.Code &&
		e.Details == t.Details &&
		e.Err == t.Err
}

// StatusCode returns the HTTP status code for the error.
func (e *CustomError) StatusCode() int {
	switch e.Code {
	case ErrWrongInput:
		return http.StatusBadRequest
	case ErrWrongCredentials:
		return http.StatusUnauthorized
	case ErrDuplicateEmail:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func NewInternalServerError(err error) *CustomError {
	return &CustomError{
		Code: ErrUnknown,
		Err:  errors.WrapPrefix(err, "Internal Server Exception", 0),
	}
}

// NewWrongCredentialsError returns a new error with the given message.
// This error should be used when the user provides wrong credentials.
func NewWrongCredentialsError(err error) error {
	return &CustomError{
		Code: ErrWrongCredentials,
		Err:  errors.WrapPrefix(err, "Wrong Credentials", 0),
	}
}

// NewWrongInputError returns a new error with the given message.
// This error should be used when wrong input is provided.
func NewWrongInputError(err error) error {
	return &CustomError{
		Code: ErrWrongInput,
		Err:  errors.WrapPrefix(err, "Wrong input", 0),
	}
}

// NewDuplicateEmailError returns an error that should be used
// when the user tries to register with an email that is already in the system.
func NewDuplicateEmailError() error {
	return &CustomError{
		Code: ErrDuplicateEmail,
		Err:  errors.New("Duplicate email"),
	}
}

// NewUnauthorizedError returns an error that should be used
// when the user is not authorized to perform the operation.
func NewUnauthorizedError(err error) error {
	return &CustomError{
		Code: ErrUnauthorized,
		Err:  errors.WrapPrefix(err, "Unauthorized error", 0),
	}
}
