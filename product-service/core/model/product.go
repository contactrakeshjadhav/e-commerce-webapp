package model

import (
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/dto"
)

type Product struct {
	ID          ID     `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (t *Product) ToDTO() *dto.Product {
	return &dto.Product{
		ID:          string(t.ID),
		Name:        t.Name,
		Description: t.Description,
		Color:       t.Color,
	}
}

func ConvertSaveProductDTOToModel(saveProductDTO dto.Product) Product {
	return Product{
		ID:          ID(saveProductDTO.ID),
		Name:        saveProductDTO.Name,
		Description: saveProductDTO.Description,
		Color:       saveProductDTO.Color,
	}
}

type Resource struct {
	// ID is the UUID of any resource
	ID string `json:"id" example:"1854ea07-01e5-42e4-b3c6-f1a5f7f150e9"`
}

func (r *Resource) Validate() error {
	if r.ID == "" {
		return dto.NewRequiredFieldError(r.ID)
	}
	return nil
}
