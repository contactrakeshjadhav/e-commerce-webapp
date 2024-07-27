package dto

type ItemsResponse struct {
	Items       interface{} `json:"items"`
	HasNext     bool        `json:"hasNext"`
	HasPrevious bool        `json:"hasPrevious"`
	TotalItems  int         `json:"totalItems"`
	TotalPages  int         `json:"totalPages"`
	PageSize    int         `json:"pageSize"`
} // @name commons.ItemsResponse

func NewItemsResponse(items interface{}, page Paging, totalItems int) ItemsResponse {
	if page.Size == 0 && page.Page == 0 {
		return ItemsResponse{
			Items:       items,
			PageSize:    totalItems,
			TotalItems:  totalItems,
			HasNext:     false,
			HasPrevious: false,
			TotalPages:  1,
		}
	}

	totalPages := totalItems / page.Size
	remainder := totalItems % page.Size
	if remainder > 0 {
		totalPages++
	}

	hasNext := page.Page < totalPages
	hasPrevious := page.Page > 1

	return ItemsResponse{
		Items:       items,
		PageSize:    page.Size,
		TotalItems:  totalItems,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
		TotalPages:  totalPages,
	}
}
