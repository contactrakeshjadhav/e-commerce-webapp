package products

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/dto"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/logger"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/repository/products"
	"github.com/google/uuid"
)

const (
	ProductRecordsLimit = 20
)

func NewProductService(productRepo products.RepositoryDB, logg logger.Logger) ProductService {
	return &productService{
		ProductRepo: productRepo,
		Log:         logg,
	}
}

type productService struct {
	ProductRepo products.RepositoryDB
	Log         logger.Logger
}

func (ts *productService) AddProduct(ctx context.Context, product model.Product) (model.Product, error) {
	logg := ts.Log.WithReqID(ctx)
	product.ID = model.ID(uuid.New().String())

	// we validate that the product name doesn't exists
	existingProduct, err := ts.ProductRepo.FindByName(ctx, product.Name)
	if !errors.Is(err, model.ErrNotFound) {
		if err == nil {
			logg.Errorf("product with name %v already exist", existingProduct.Name)
			return model.Product{}, errors.Wrapf(model.ErrAlreadyExists, "product name already exist")
		}
		logg.Errorf("failed to find product by name: %v", err)
		return model.Product{}, err
	}

	newProduct, err := ts.ProductRepo.Insert(ctx, product)
	if err != nil {
		logg.Errorf("failed to create new product: %v", err)
		return model.Product{}, err
	}
	return newProduct, nil
}

func (ts *productService) UpdateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	var err error
	logg := ts.Log.WithReqID(ctx)

	// check if product already exists
	existingProduct, err := ts.ProductRepo.FindByID(ctx, product.ID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			logg.Errorf("product with id %v doesn't exist", existingProduct.ID)
			return model.Product{}, errors.Wrapf(model.ErrNotFound, "product doesn't exist")
		}
		logg.Errorf("failed to get product by id: %s", err)
		return model.Product{}, err
	}

	err = ts.ProductRepo.Update(ctx, product)
	if err != nil {
		if strings.Contains(err.Error(), model.ErrUniqueConstraint.Error()) {
			logg.Errorf("product with name %v already exist", existingProduct.Name)
			return model.Product{}, errors.Wrapf(model.ErrAlreadyExists, "product name already exist")
		}
		logg.Errorf("failed to update product: %s", err)
		return model.Product{}, err
	}

	return product, nil
}

func (ts *productService) DeleteProduct(ctx context.Context, productID string, forceDelete bool) error {
	logg := ts.Log.WithReqID(ctx)

	if !forceDelete {
		// we validate that the product is unused

		err := ts.ProductRepo.Delete(ctx, productID)
		if err != nil {
			logg.Errorf("failed to delete product: %v", err)
			return err
		}
	} else {
		// delete objects to product mapping
		err := ts.ProductRepo.Delete(ctx, productID)
		if err != nil {
			logg.Errorf("failed to delete product: %v", err)
			return err
		}
	}
	return nil
}

func (ts *productService) FindAll(ctx context.Context, pageFilterSort commons.PageFilterSort) ([]model.Product, int, error) {
	products, total, err := ts.ProductRepo.FindAll(ctx, pageFilterSort)
	if err != nil {
		ts.Log.Errorf("failed to find all products: %v", err)
		return []model.Product{}, 0, err
	}
	return products, total, nil
}

func (ts *productService) FindProductsByPartialName(ctx context.Context, name string) ([]model.Product, error) {
	product, err := ts.ProductRepo.FindByPartialName(ctx, name, ProductRecordsLimit)
	if err != nil {
		ts.Log.Errorf("failed to find all products: %v", err)
		return []model.Product{}, err
	}
	return product, nil
}

func (ts *productService) FindByID(ctx context.Context, id model.ID) (model.Product, error) {
	product, err := ts.ProductRepo.FindByID(ctx, id)
	if err != nil {
		ts.Log.Errorf("failed to find product by id: %v", err)
		return model.Product{}, err
	}
	return product, nil
}

func (ts *productService) GetProductsByIDs(ctx context.Context, ids []string) ([]model.Product, error) {
	products, err := ts.ProductRepo.GetProductsByIDs(ctx, ids)
	if err != nil {
		ts.Log.Errorf("failed to get products by ids: %v", err)
		return []model.Product{}, err
	}
	return products, nil
}
