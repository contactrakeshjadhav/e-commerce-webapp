package controllers

import (
	"net/http"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/swagger"
)

// SwaggerDocsHandler loads and expose our openapiV2 definition in json format
func SwaggerDocsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type:", "application/json")
	w.Write(swagger.JsonFile)
}
