package model

import (
	"errors"
	"fmt"
)

// General errors for product API
var (
	ErrEnqueueJobFailed        = errors.New("error creating job")
	ErrJsonMarshallingFailed   = errors.New("error marshalling the desired object")
	ErrJsonUnmarshallingFailed = errors.New("error unmarshalling the desired object")
	ErrNotFound                = errors.New("not found")
	ErrAlreadyExists           = errors.New("already exists")
	ErrInternalServerFail      = errors.New("internal server error")
	ErrInvalidArgument         = errors.New("invalid argument")
	ErrPermissionDenied        = errors.New("permission denied")
	ErrUserInfoValidation      = errors.New("failed to validate user information")
	ErrQueryFailed             = errors.New("error query failed")
	ErrBadRequest              = errors.New("body of request contains errors")
	ErrUniqueConstraint        = errors.New("duplicate key value violates unique constraint")
)

func ErrUnexpectedHTTPCode(want, got int) error {
	return fmt.Errorf("unexpected HTTP status code, want %v - got %v", want, got)
}

func ErrJsonDecodeInputMessage(function string, msg string) string {
	return fmt.Sprintf("failed to decode %v input: %v", function, msg)
}

func ErrJsonEncodeResponseMessage(function string, msg string) string {
	return fmt.Sprintf("failed to decode %v response: %v", function, msg)
}

func ErrFailedToValidateInputMessage(function string, msg string) string {
	return fmt.Sprintf("failed to validate %v input: %v", function, msg)
}
