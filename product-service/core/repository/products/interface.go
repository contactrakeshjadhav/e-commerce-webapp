package products

import (
	"context"

	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/dto"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
)

type RepositoryDB interface {
	// reader
	FindByID(ctx context.Context, id model.ID) (model.Product, error)
	FindByName(ctx context.Context, name string) (model.Product, error)
	FindAll(ctx context.Context, param commons.PageFilterSort) ([]model.Product, int, error)
	FindByPartialName(ctx context.Context, name string, limit int32) ([]model.Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]model.Product, error)

	// writer
	Insert(ctx context.Context, product model.Product) (model.Product, error)
	Update(ctx context.Context, tool model.Product) error
	Delete(ctx context.Context, id string) error
}
