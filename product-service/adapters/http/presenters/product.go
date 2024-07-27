package presenters

import (
	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/dto"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/dto"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
)

// ProductItem this is a presenter responsible to format a model.Product
// return a single dto.Product{}
func ProductItem(product model.Product) dto.Product {
	return *product.ToDTO()
}

// ProductCollection this is a presenter responsible to format a collection of model.Product
// return a list []dto.Products{}
func ProductCollection(products []model.Product) []dto.Product {
	productCollection := make([]dto.Product, len(products))
	for i, product := range products {
		productCollection[i] = ProductItem(product)
	}
	return productCollection
}

// ObjectProductsCollection this is a presenter responsible to format a collection of model.Product
// return a list map[string][]dto.Products{}
func ObjectProductsCollection(objectProducts map[string][]model.Product) map[string][]dto.Product {
	objectProductsDto := map[string][]dto.Product{}
	for k, products := range objectProducts {
		objectProductsDto[k] = ProductCollection(products)
	}
	return objectProductsDto
}

func NewProductsResponse(items interface{}, page commons.Paging, totalItems int) commons.ItemsResponse {
	return commons.NewItemsResponse(items, page, totalItems)
}
