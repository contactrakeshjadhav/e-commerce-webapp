package dto

import (
	"fmt"
	"regexp"
	"strings"

	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/dto"
)

type Product struct {
	// ID is UUID of the product to update, it's empty when adding a new product. [optional value ? empty string]
	ID string `json:"id" example:"26d661b4-306c-43c5-a57b-db8c1e813f72"`
	// Name is product name.
	Name string `json:"name" example:"SAS Programs"`
	// Description is product description. [optional value ? empty string]
	Description string `json:"description" example:"This is SAS Programs product"`
	// Color is product color.
	Color string `json:"color" example:"#126"`
} // @Name Product

type DeleteProduct struct {
	// Id is the UUID of the product to be deleted.
	ID string `json:"id" example:"d0fb502e-5ed3-4c7b-8e1c-54cd813378e6"`
	// ForceDelete is used to delete used products aswell when set to true. Default is false.
	ForceDelete bool `json:"forceDelete" example:"false"`
} // @Name DeleteProduct

type GetProducts struct {
	// optional parameter for sorting/paging/filtering, if empty, we return all the records
	PageFilterSort commons.PageFilterSort `json:"pageFilterSort"`
} // @Name GetProducts

func (product *Product) Validate() error {

	if strings.TrimSpace(string(product.Name)) == "" {
		return NewRequiredFieldError("name")

	}

	if strings.TrimSpace(string(product.Color)) == "" {
		return NewRequiredFieldError("color")

	}

	if !isValidHexColorCode(product.Color) {
		return ErrInvalidColorCode
	}
	product.Color = checkPrefixHash(product.Color)
	return nil
}

type GetProductsByName struct {
	// Name is product name.
	Name string `json:"name" example:"SAS Programs"`
} //@Name GetProductsByName

type GetProductsByIDs struct {
	// IDs are the product UUIDs.
	IDs []string `json:"ids" example:"ff81c7b4-b4d5-4708-aa99-cbcefb37cb0d, 022aa9e2-c7cb-4a03-9fa2-6ef034f4a2b2"`
} //@Name GetProductsByIDs

type GetProductsByObjectIDs struct {
	// ObjectIDs are array of object UUIDs.
	ObjectIDs []string `json:"objectIDs" example:"170f8cdf-f71e-4c57-9d1a-2976879fb21a, 022aa9e2-c7cb-4a03-9fa2-6ef034f4a2b2"`
} //Name GetProductsByObjectIDs

func (product *DeleteProduct) Validate() error {
	if strings.TrimSpace(product.ID) == "" {
		return NewRequiredFieldError("id")
	}
	return nil
}

func (products *GetProductsByIDs) Validate() error {
	if len(products.IDs) == 0 {
		return NewRequiredFieldError("ids")
	}
	return nil
}

func (objects *GetProductsByObjectIDs) Validate() error {
	if len(objects.ObjectIDs) == 0 {
		return NewRequiredFieldError("objectIDs")
	}
	return nil
}

// Check "#" prefix exists or not...
func checkPrefixHash(colorCode string) string {
	if len(colorCode) > 0 && string(colorCode[0]) != "#" {
		colorCode = fmt.Sprintf("#%s", colorCode)
	}
	return colorCode
}

func isValidHexColorCode(colorCode string) bool {
	// Define a regular expression pattern for a valid hex color code
	// The pattern matches codes with or without the '#' prefix and of length 3 or 6 characters.
	pattern := regexp.MustCompile(`^#?([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)

	// Use the pattern to match the input hex color code
	return pattern.MatchString(colorCode)
}

var (
	GetProductsNameColumn        string = "name"
	GetProductsDescriptionColumn string = "description"
	GetProductsColorColumn       string = "color"

	GetProductsFilterColumns = map[string]struct{}{
		GetProductsNameColumn:        {},
		GetProductsDescriptionColumn: {},
		GetProductsColorColumn:       {},
	}
	GetProductsSortColumns = map[string]struct{}{
		GetProductsNameColumn:        {},
		GetProductsDescriptionColumn: {},
		GetProductsColorColumn:       {},
	}
)
