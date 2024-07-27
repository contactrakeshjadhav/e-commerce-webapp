package dto

import (
	"errors"
	"fmt"
)

var (
	ErrRequiredFieldIsEmpty = errors.New("required field is empty")
	ErrInvalidColorCode     = errors.New("invalid color code")
	ErrEmptyproductID       = errors.New("product id can't be empty")
	ErrEmptyName            = errors.New("name can't be empty")
)

func NewRequiredFieldError(field string) error {
	return fmt.Errorf("field %v is required", field)
}

func InvalidValueError(field, value string) error {
	return fmt.Errorf("invalid value '%v' for '%v'", value, field)
}
