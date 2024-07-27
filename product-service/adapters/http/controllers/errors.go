package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	ErrUnexpected error = errors.New("an error occurred while processing your request")
)

type Error struct {
	Message string `json:"message"`
} //@name Error

// @Description ApplicationError is a generic structure that represents an error response.
type ApplicationError struct {
	Error Error `json:"error"`
} //@name ApplicationError

func InvalidArgumentError(argument interface{}, message string) error {
	return fmt.Errorf("invalid argument %v, msg = %v", argument, message)
}

func CheckErrors(err error) string {
	var unmarshalErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError
	var errorMessage string

	switch {
	case errors.As(err, &unmarshalErr):
		errorMessage = fmt.Sprintf(
			"%s is a %s but %s was expected",
			unmarshalErr.Field,
			unmarshalErr.Value,
			unmarshalErr.Type.Kind().String(),
		)
	case errors.As(err, &syntaxErr):
		errorMessage = fmt.Sprintf("syntax error at the bytes offset %d", syntaxErr.Offset)
	case errors.Is(err, io.EOF):
		errorMessage = "request body is missing"
	default:
		errorMessage = "input body is in a bad format"
	}
	return errorMessage
}
