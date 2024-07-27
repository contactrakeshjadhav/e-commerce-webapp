package dto

import "math"

const DEFAULT_PAGE = 1
const DEFAULT_SIZE = math.MaxInt32

type PageFilterSort struct {
	Paging Paging `json:"paging"`
	Sort   Sort   `json:"sort"`
	Filter Filter `json:"filter"`
} // @name commons.PageFilterSort
