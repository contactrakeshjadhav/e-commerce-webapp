package products

import (
	"context"

	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/dto"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
)

type ProductService interface {
	AddProduct(ctx context.Context, product model.Product) (model.Product, error)
	UpdateProduct(ctx context.Context, tool model.Product) (model.Product, error)
	FindProductsByPartialName(ctx context.Context, name string) ([]model.Product, error)
	DeleteProduct(ctx context.Context, productID string, forceDelete bool) error
	FindAll(ctx context.Context, pageFilterSort commons.PageFilterSort) ([]model.Product, int, error)
	FindByID(ctx context.Context, id model.ID) (model.Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]model.Product, error)
}
