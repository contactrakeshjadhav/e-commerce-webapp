package products

import (
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
	productDatabase "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/database/sqlc-generated/queries"
)

func entityproductToModel(entity productDatabase.Product) model.Product {
	return model.Product{
		ID:          model.ID(entity.ID),
		Name:        entity.Name.String,
		Description: entity.Description.String,
		Color:       entity.Color.String,
	}
}

func productEntityToModel(collection []productDatabase.Product) []model.Product {
	products := make([]model.Product, len(collection))
	for i, j := range collection {
		j := entityproductToModel(j)
		products[i] = j
	}
	return products
}
func findProductsManyEntitiesToModel(entities []productDatabase.FindManyProductsRow) (products []model.Product, total int64) {
	if len(entities) == 0 {
		return
	}
	products = make([]model.Product, len(entities))
	for i, entity := range entities {
		products[i] = model.Product{
			ID:          model.ID(entity.ID),
			Name:        entity.Name.String,
			Description: entity.Description.String,
			Color:       entity.Color.String,
		}
	}
	return products, entities[0].Count
}
