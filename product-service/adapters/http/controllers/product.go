package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/adapters/http/presenters"
	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/model"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/logger"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/dto"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
	product "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/services/products"
)

const (
	saveProductMethod              string = "saveProduct"
	deleteProductMethod            string = "deleteProduct"
	getProductsMethod              string = "getProducts"
	getProductsByPartialNameMethod string = "getProductsByPartialName"
	getProductMethod               string = "getProduct"
	getProductsByIDsMethod         string = "getProductsByIDs"
	getProductsByObjectIDsMethod   string = "getProductsByObjectIDs"
)

type ProductController struct {
	ProductService product.ProductService
	l              logger.Logger
}

func NewProductController(productService product.ProductService, l logger.Logger) *ProductController {
	return &ProductController{
		ProductService: productService,
		l:              l,
	}
}

// save product
// @Summary      add or update a product
// @Description  This API adds a new product if id field is empty. If it's not empty, it tries to update product with that id.
// @Products         products
// @Accept       json
// @Param        payload   body dto.Product true "payload"
// @Produce      json
// @Success      201  {object} dto.Product
// @Success      200  {object} dto.Product
// @Failure      400  {object} ApplicationError
// @Failure      500  {object} ApplicationError
// @Failure      403  {object} ApplicationError
// @Router       /v1/saveProduct [post]
func (tc *ProductController) SaveProduct(w http.ResponseWriter, r *http.Request) {
	log := tc.l.WithReqID(r.Context())
	var saveProductDTO dto.Product
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&saveProductDTO)
	if err != nil {
		errMessage := CheckErrors(err)
		log.Errorf(model.ErrJsonDecodeInputMessage(saveProductMethod, errMessage))
		ApplicationErrorResponse(w, errors.New(errMessage), http.StatusBadRequest)
		return
	}
	log.LogDTO(saveProductDTO)

	if err := saveProductDTO.Validate(); err != nil {
		log.Infof(model.ErrFailedToValidateInputMessage(saveProductMethod, err.Error()))
		ApplicationErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	product := model.ConvertSaveProductDTOToModel(saveProductDTO)
	var httpStatus int
	var action string
	if product.ID != "" {
		product, err = tc.ProductService.UpdateProduct(r.Context(), product)
		httpStatus = http.StatusOK
		action = "update"
	} else {
		product, err = tc.ProductService.AddProduct(r.Context(), product)
		httpStatus = http.StatusCreated
		action = "add"
	}

	if err != nil {
		log.Errorf("failed to %v product: %v", action, err)
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	output := presenters.ProductItem(product)
	result, err := json.Marshal(output)
	if err != nil {
		log.Errorf(model.ErrJsonEncodeResponseMessage(saveProductMethod, err.Error()))
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	log.LogDTO(output)
	JsonResponse(w, result, httpStatus)
}

// delete product
// @Summary      delete an unused product
// @Description  This API deletes an unused product if force_delete is False. If force_delete is True then it deletes used product as well.
// @Products         products
// @Accept       json
// @Param        payload   body dto.DeleteProduct true "payload"
// @Produce      json
// @Success      200  {object} dto.DeleteProduct
// @Failure      400  {object} ApplicationError
// @Failure      500  {object} ApplicationError
// @Failure      403  {object} ApplicationError
// @Router       /v1/deleteProduct [post]
func (tc *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	log := tc.l.WithReqID(r.Context())
	var deleteProductDTO dto.DeleteProduct
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&deleteProductDTO)
	if err != nil {
		errMessage := CheckErrors(err)
		log.Errorf(model.ErrJsonDecodeInputMessage(deleteProductMethod, errMessage))
		ApplicationErrorResponse(w, errors.New(errMessage), http.StatusBadRequest)
		return
	}
	log.LogDTO(deleteProductDTO)

	if err := deleteProductDTO.Validate(); err != nil {
		log.Infof(model.ErrFailedToValidateInputMessage(deleteProductMethod, err.Error()))
		ApplicationErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = tc.ProductService.DeleteProduct(r.Context(), deleteProductDTO.ID, deleteProductDTO.ForceDelete)
	if err != nil {
		log.Errorf("failed to delete product: %v", err)
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	JsonResponse(w, nil, http.StatusOK)
}

// list products
// @Summary      list products
// @Description  list all available products
// @Description  * Can be sorted by: name, description, color
// @Description  * Can be filtered by: name
// @Param        payload	body dto.GetProducts true  "A wrapper for the paging, sorting and filtering object."
// @Products         products
// @Accept       json
// @Produce      json
// @Success      200  {object} commons.ItemsResponse{items=[]dto.Product}
// @Failure      400  {object} ApplicationError
// @Failure      500  {object} ApplicationError
// @Router       /v1/getProducts [post]
func (tc *ProductController) GetProducts(rw http.ResponseWriter, r *http.Request) {
	logg := tc.l.WithReqID(r.Context())
	input := dto.GetProducts{}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errMessage := CheckErrors(err)
		logg.Errorf("failed to decode request input: %v", errMessage)
		ApplicationErrorResponse(rw, errors.New(errMessage), http.StatusBadRequest)
		return
	}
	logg.LogDTO(input)

	pageFilterSort := input.PageFilterSort
	pageFilterSort.Paging.FillDefaults()

	err = pageFilterSort.Filter.Validate(dto.GetProductsFilterColumns)
	if err != nil {
		ApplicationErrorResponse(rw, err, http.StatusBadRequest)
		return
	}
	err = pageFilterSort.Sort.Validate(dto.GetProductsSortColumns)
	if err != nil {
		ApplicationErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	// Get products
	products, total, err := tc.ProductService.FindAll(r.Context(), pageFilterSort)
	if err != nil {
		logg.Errorf("failed to find all products: %v", err)
		ApplicationErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//prepare presentation
	productsCollection := presenters.ProductCollection(products)
	items := presenters.NewProductsResponse(productsCollection, pageFilterSort.Paging, total)
	//format
	result, err := json.Marshal(items)
	if err != nil {
		logg.Errorf(model.ErrJsonEncodeResponseMessage(getProductsMethod, err.Error()))
		ApplicationErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}
	logg.LogDTO(items)
	JsonResponse(rw, result, http.StatusOK)
}

// list products
// @Summary      list products by name
// @Description  list all available products by name
// @Description  * Can be filtered by: name
// @Param        payload	body dto.GetProductsByName true "payload"
// @Products         products
// @Accept       json
// @Produce      json
// @Success      200  {array} dto.Product
// @Failure      400  {object} ApplicationError
// @Failure      500  {object} ApplicationError
// @Router       /v1/getProductsByPartialName [post]
func (tc *ProductController) GetProductsByPartialName(rw http.ResponseWriter, r *http.Request) {
	logg := tc.l.WithReqID(r.Context())
	input := dto.GetProductsByName{}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errMessage := CheckErrors(err)
		logg.Errorf("failed to decode request input: %v", errMessage)
		ApplicationErrorResponse(rw, errors.New(errMessage), http.StatusBadRequest)
		return
	}
	logg.LogDTO(input)

	products, err := tc.ProductService.FindProductsByPartialName(r.Context(), input.Name)
	if err != nil {
		logg.Errorf("failed to find all product name: %v", err)
		ApplicationErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}
	output := presenters.ProductCollection(products)
	result, err := json.Marshal(output)
	if err != nil {
		logg.Errorf(model.ErrJsonEncodeResponseMessage(getProductsByPartialNameMethod, err.Error()))
		ApplicationErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}
	logg.LogDTO(output)
	JsonResponse(rw, result, http.StatusOK)
}

// @Summary      get product
// @Description  Returns product specification by id.
// @Products         products
// @Param 	     id	body model.Resource true  "Product id"
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.Product
// @Failure      400  {object}  ApplicationError
// @Failure      500  {object}  ApplicationError
// @Router       /v1/getProduct [post]
func (tc *ProductController) GetProduct(rw http.ResponseWriter, r *http.Request) {
	logg := tc.l.WithReqID(r.Context())
	input := model.Resource{}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errMessage := CheckErrors(err)
		logg.
			WithInt(commons.ResponseCodeKey.String(), http.StatusBadRequest).
			Errorf("failed to decode request input: %v", errMessage)
		ApplicationErrorResponse(rw, errors.New(errMessage), http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		ApplicationErrorResponse(rw, err, http.StatusBadRequest)
		return
	}
	logg.LogDTO(input)

	product, err := tc.ProductService.FindByID(r.Context(), model.ID(input.ID))
	if err != nil {
		logg.Errorf("failed to find product by id: %v", err)
		ApplicationErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}
	output := presenters.ProductItem(product)
	result, err := json.Marshal(output)
	if err != nil {
		logg.Errorf(model.ErrJsonEncodeResponseMessage(getProductMethod, err.Error()))
		ApplicationErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}
	logg.LogDTO(output)
	JsonResponse(rw, result, http.StatusOK)
}

// @Summary      get products by ids
// @Description  Returns array of products specification by given ids.
// @Products         products
// @Param        payload  body dto.GetProductsByIDs true "payload"
// @Accept       json
// @Produce      json
// @Success      200  {object}  []dto.Product
// @Failure      400  {object}  ApplicationError
// @Failure      500  {object}  ApplicationError
// @Router       /v1/getProductsByIDs [post]
func (tc *ProductController) GetProductsByIDs(rw http.ResponseWriter, r *http.Request) {
	logg := tc.l.WithReqID(r.Context())
	input := dto.GetProductsByIDs{}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errMessage := CheckErrors(err)
		logg.
			WithInt(commons.ResponseCodeKey.String(), http.StatusBadRequest).
			Errorf("failed to decode request input: %v", errMessage)
		ApplicationErrorResponse(rw, errors.New(errMessage), http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		ApplicationErrorResponse(rw, err, http.StatusBadRequest)
		return
	}
	logg.LogDTO(input)

	products, err := tc.ProductService.GetProductsByIDs(r.Context(), input.IDs)
	if err != nil {
		logg.Errorf("failed to get products by ids: %v", err)
		ApplicationErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}
	output := presenters.ProductCollection(products)
	result, err := json.Marshal(output)
	if err != nil {
		logg.Errorf(model.ErrJsonEncodeResponseMessage(getProductsByIDsMethod, err.Error()))
		ApplicationErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}
	logg.LogDTO(output)
	JsonResponse(rw, result, http.StatusOK)
}
