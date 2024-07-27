package dto

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidArchivableValue error = errors.New("values for archivable filter can be: true/false")
)

func NewRequiredFieldError(field string) error {
	return fmt.Errorf("field %v is required", field)
}
