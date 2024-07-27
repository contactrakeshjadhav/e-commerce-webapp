package dto

import (
	"errors"
	"strings"
)

var (
	ErrSortByFieldNotFound error = errors.New("can't sort by field")
	ErrUnknownSortDir      error = errors.New("unknown sort direction")
)

type Sort struct {
	// SortBy Name of the property which need to be sorted
	SortBy string `json:"sortBy" example:"name"`
	// SortDir The direction of the sorting: ASC or DESC. Default = DESC
	SortDir string `json:"sortDir" swaggerType:"string" enums:"ASC,DESC"`
} // @name commons.Sort

// Validate Will validate input for sorting
func (s *Sort) Validate(expected map[string]struct{}) error {
	if !s.Present() {
		return nil
	}

	if _, ok := expected[s.SortBy]; !ok {
		return ErrSortByFieldNotFound
	}

	dir := GetSortDirectionFromString(strings.ToUpper(s.SortDir))
	if !dir.IsValid() {
		return ErrUnknownSortDir
	}

	return nil
}

func (s *Sort) Present() bool {
	return s.SortBy != ""
}
