package products

import (
	"context"
	"database/sql"
	"fmt"

	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/dto"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/logger"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
	database "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/database/sqlc-generated/queries"
	"github.com/pkg/errors"
)

var (
	ErrInvalidPaging = errors.New("page number and size must be greater than 0")
)

const (
	nameColumn = "name"
)

type productRepository struct {
	DB      *sql.DB
	queries *database.Queries
	logger  logger.Logger
}

func NewProductRepository(db *sql.DB, logg logger.Logger) RepositoryDB {
	return &productRepository{
		DB:      db,
		queries: database.New(db),
		logger:  logg,
	}
}

func (tr *productRepository) Insert(ctx context.Context, product model.Product) (model.Product, error) {
	log := tr.logger.WithReqID(ctx)
	tx, err := tr.DB.Begin()
	if err != nil {
		log.Errorf("failed to start transaction for product insertion: %v", err)
		return model.Product{}, errors.Wrap(err, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()

	q := tr.queries.WithTx(tx)
	err = q.InsertProduct(ctx, database.InsertProductParams{
		ID:          string(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Color:       product.Color,
	})

	if err != nil {
		log.Errorf("failed to insert product: %v", err)
		return model.Product{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commit product insertion: %v", err)
		return model.Product{}, errors.New("failed to commit insertion of a new product")
	}

	return product, nil
}

func (tr *productRepository) Update(ctx context.Context, product model.Product) error {
	log := tr.logger.WithReqID(ctx)
	tx, err := tr.DB.Begin()
	if err != nil {
		log.Errorf("failed to start transaction for product update: %v", err)
		return errors.Wrap(err, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()

	q := tr.queries.WithTx(tx)

	err = q.UpdateProduct(ctx, database.UpdateProductParams{
		Name:        product.Name,
		ID:          string(product.ID),
		Description: product.Description,
		Color:       product.Color,
	})

	if err != nil {
		log.Errorf("failed to update product: %v", err)
		return err
	}

	return tx.Commit()
}

// FindByID return the product related to the desire ID
// return an error on any operation fail
func (tr *productRepository) FindByID(ctx context.Context, id model.ID) (model.Product, error) {
	logg := tr.logger.WithReqID(ctx)

	productEntity, err := tr.queries.FindProductById(ctx, string(id))
	if err != nil {
		errMessage := fmt.Sprintf("failed to find product by id (%v): %v", string(id), err)
		if errors.Is(err, sql.ErrNoRows) {
			logg.Infof(errMessage)
			return model.Product{}, errors.Wrapf(model.ErrNotFound, err.Error())
		}
		logg.Errorf(errMessage)
		return model.Product{}, errors.Wrapf(model.ErrQueryFailed, err.Error())
	}

	product := entityproductToModel(productEntity)
	return product, nil
}

// FindByName return the product related to the given name
// return an error on any operation fail
func (tr *productRepository) FindByName(ctx context.Context, name string) (model.Product, error) {
	logg := tr.logger.WithReqID(ctx)

	productEntity, err := tr.queries.FindProductByName(ctx, name)
	if err != nil {
		errMessage := fmt.Sprintf("failed to find product by name (%v): %v", name, err)
		if errors.Is(err, sql.ErrNoRows) {
			logg.Infof(errMessage)
			return model.Product{}, errors.Wrapf(model.ErrNotFound, err.Error())
		}
		logg.Errorf(errMessage)
		return model.Product{}, errors.Wrapf(model.ErrQueryFailed, err.Error())
	}

	product := entityproductToModel(productEntity)
	return product, nil
}

// This function will delete the product entity for given id
func (tr *productRepository) Delete(ctx context.Context, productID string) error {
	logg := tr.logger.WithReqID(ctx)

	err := tr.queries.DeleteProductById(ctx, productID)
	if err != nil {
		logg.Errorf("failed to delete product by id %v: %v", productID, err)
		return errors.Wrap(model.ErrInternalServerFail, err.Error())
	}
	return nil
}

const DefaultproductsOrderByColumn = "name"

func (tr *productRepository) FindAll(ctx context.Context, param commons.PageFilterSort) ([]model.Product, int, error) {
	return tr.findManyproducts(ctx, nil, param)
}

// findMany returns all the records according to the paging, sorting and filtering.
func (tr *productRepository) findManyproducts(ctx context.Context, ids []model.ID, param commons.PageFilterSort) ([]model.Product, int, error) {
	var orderBy string
	if !param.Sort.Present() {
		orderBy = DefaultproductsOrderByColumn
		param.Sort.SortDir = commons.SORT_ASC.String()
	} else {
		orderBy = param.Sort.SortBy
	}

	if param.Paging.Size <= 0 || param.Paging.Page <= 0 {
		return []model.Product{}, 0, ErrInvalidPaging
	}

	offset := getOffset(param.Paging)
	productEntities, err := tr.queries.FindManyProducts(ctx, database.FindManyProductsParams{
		PageOffset: offset,
		PageLimit:  int32(param.Paging.Size),
		OrderBy:    orderBy,
		Ascending:  param.Sort.SortDir == commons.SORT_ASC.String(),
		FilterBy:   param.Filter.FilterBy,
		Pattern:    "%" + param.Filter.FilterPattern + "%",
		Ids:        IDsToStrings(ids),
	})
	if err != nil {
		tr.logger.WithReqID(ctx).Errorf("Could not find products: %v", err)
		return []model.Product{}, 0, err
	}
	products, total := findProductsManyEntitiesToModel(productEntities)
	var productPointers = make([]*model.Product, len(products))
	for i, product := range products {
		productPointers[i] = &product
	}
	return products, int(total), nil
}

// FindByPartialName return the product related to the given name
// return an error on any operation fail
func (tr *productRepository) FindByPartialName(ctx context.Context, name string, productRecordsLimit int32) ([]model.Product, error) {
	logg := tr.logger.WithReqID(ctx)

	productEntity, err := tr.queries.FindProductsByPartialName(ctx, database.FindProductsByPartialNameParams{
		OrderBy:   DefaultproductsOrderByColumn,
		FilterBy:  nameColumn,
		Ascending: true,
		PageLimit: productRecordsLimit,
		Pattern:   name + "%",
	})
	if err != nil {
		logg.Infof("failed to find product by name: %s", name)
		return []model.Product{}, errors.Wrap(model.ErrQueryFailed, err.Error())
	}
	product := productEntityToModel(productEntity)
	return product, nil
}

// GetproductsByIDs return the products related to the given IDs
// return an error on any operation fail
func (tr *productRepository) GetProductsByIDs(ctx context.Context, ids []string) ([]model.Product, error) {
	logg := tr.logger.WithReqID(ctx)

	productEntities, err := tr.queries.GetProductsByIDs(ctx, ids)
	if err != nil {
		errMessage := fmt.Sprintf("failed to get products by ids (%v): %v", ids, err)
		logg.Errorf(errMessage)
		return []model.Product{}, errors.Wrapf(model.ErrQueryFailed, err.Error())
	}

	products := productEntityToModel(productEntities)
	return products, nil
}

func getOffset(pagination commons.Paging) int32 {
	offset := (pagination.Page - 1) * pagination.Size
	if offset < 0 {
		offset = 0
	}
	return int32(offset)
}

func IDsToStrings(ids []model.ID) []string {
	strings := make([]string, len(ids))
	for i, id := range ids {
		strings[i] = string(id)
	}
	return strings
}
