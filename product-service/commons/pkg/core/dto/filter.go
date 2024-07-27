package dto

import (
	"errors"
)

var (
	ErrFilterByFieldNotFound error = errors.New("can't filter by field")
)

type Filter struct {
	// FilterBy Name of the property which need to be filtered
	FilterBy string `json:"filterBy" example:"name"`
	// FilterPattern Search for this pattern in the column specified by filterBy. Regular expressions are not allowed
	FilterPattern string `json:"filterPattern" example:"project"`
} // @name commons.Filter

// Validate Will validate fiter pattern. Regular expressions are not allowed
func (f *Filter) Validate(expected map[string]struct{}) error {
	if !f.Present() {
		return nil
	}

	if _, ok := expected[f.FilterBy]; !ok {
		return ErrFilterByFieldNotFound
	}

	// @TODO https://d-wise.atlassian.net/browse/ASPL-875
	return nil
}

func (f *Filter) Present() bool {
	return f.FilterBy != ""
}
