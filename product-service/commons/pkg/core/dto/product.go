package dto

import (
	"strings"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/utils"
)

type product struct {
	// ID is ID of the product to update, it's empty when adding a new product. [optional value ? empty string]
	ID string `json:"id"`
	// Name is product name.
	Name string `json:"name"`
	// Description is product description. [optional value ? empty string]
	Description string `json:"description"`
	// Color is product color.
	Color string `json:"color"`
} // @Name product

type productsToObject struct {
	// productIDs is a list of product IDs
	productIDs []string `json:"productIDs"`
	// ObjectID is the ID of the object that needs to be mapped with products.
	ObjectID string `json:"objectId"`
	// ObjectType is the type of object that is mapped with products.
	ObjectType string `json:"objectType"`
} // @Name productsToObject

type FindObjectsByproductsInput struct {
	// productIds list of product IDs
	productIds []string `json:"productIDs"`
	// Object type
	ObjectType string `json:"objectType"`
	// Operator
	Operator string `json:"operator" enums:"OR,AND"`
} // @Name FindObjectsByproductsInput

type GetproductsByObjectIDs struct {
	// ObjectIDs are array of object ids.
	ObjectIDs []string `json:"objectIDs"`
} //@Name GetproductsByObjectIDs

func (item *FindObjectsByproductsInput) Validate() error {
	if len(item.productIds) == 0 {
		return NewRequiredFieldError("productIDs")
	}
	if strings.TrimSpace(item.ObjectType) == utils.EMPTY_SPACE {
		return NewRequiredFieldError("objectType")
	}
	if strings.TrimSpace(item.Operator) == utils.EMPTY_SPACE {
		return NewRequiredFieldError("operator")
	}
	return nil
}
