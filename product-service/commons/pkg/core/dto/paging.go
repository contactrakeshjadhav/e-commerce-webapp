package dto

import "database/sql"

type Paging struct {
	// Page Page Number
	Page int `json:"page" example:"1"`
	// Size Number of items to return in the response
	Size int `json:"size" example:"10"`
} // @name commons.Paging

// NewPaging Will initialize Paging struct with given page and size values.
// If those are not valid, will initialize with defaults
func NewPaging(page int, size int) Paging {
	var paging Paging
	if page <= 0 {
		paging.Page = DEFAULT_PAGE
	} else {
		paging.Page = page
	}

	if size <= 0 || size > DEFAULT_SIZE {
		paging.Size = DEFAULT_SIZE
	} else {
		paging.Size = size
	}

	return paging
}

// FillDefaults Will fill with defaults if needed
func (p *Paging) FillDefaults() {
	if p.Page <= 0 {
		p.Page = DEFAULT_PAGE
	}

	if p.Size <= 0 || p.Size > DEFAULT_SIZE {
		p.Size = DEFAULT_SIZE
	}
}

// For use with the SQL OFFSET clause
func (p *Paging) GetOffset() int32 {
	offset := (p.Page - 1) * p.Size
	if offset < 0 {
		offset = 0
	}
	return int32(offset)
}

// For use with the SQL LIMIT clause
func (p *Paging) GetLimit() sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(p.Size),
		Valid: p.Size > 0,
	}
}
