package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/model"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/logger"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
)

const (
	somethingWentWrongMessage string = "something went wrong, please try again later"
)

type BaseController struct {
	BuildInfo model.BuildInfo
	log       logger.Logger
}

func NewBaseController(info model.BuildInfo, l logger.Logger) *BaseController {
	return &BaseController{
		BuildInfo: info,
		log:       l,
	}
}

// Index - Returns all the available APIs
// @Summary This API can be used as health check for this application.
// @Description Tells if the chi-swagger APIs are working or not.
// @Tags chi-swagger
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router / [get]
func (c *BaseController) Index(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Aspire tag service - version:%v", c.BuildInfo.Version)
	message := []byte(msg)
	JsonResponse(w, message, http.StatusOK)
}

// GetBuildInfo - Returns the service build information
// Summary This API will summarize the current build information of the service.
// Description Resume the service build information.
// Tags build internal-only
// Accept  json
// Produce  json
// Success 200 {object} model.BuildInfo
// Router /getBuildInfo [post]
func (c *BaseController) GetBuildInfo(w http.ResponseWriter, r *http.Request) {
	info := c.BuildInfo
	marshaled, err := json.Marshal(&info)
	if err != nil {
		errorResponse, err := json.Marshal(ApplicationError{Error: Error{somethingWentWrongMessage}})
		if err != nil {
			c.log.WithReqID(r.Context()).
				WithInt(commons.ResponseCodeKey.String(), http.StatusInternalServerError).
				Fatalf("failed to marshal error response: %v", err)
		}
		JsonResponse(w, errorResponse, http.StatusInternalServerError)
		return
	}
	JsonResponse(w, marshaled, http.StatusOK)
}

// Ping
// @Summary This API can be used as health check for this application.
// @Description Tells if the chi-swagger APIs are working or not.
// @Tags chi-swagger
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router /ping [get]
func (c *BaseController) Ping(w http.ResponseWriter, r *http.Request) {
	message := []byte("Server is up!")
	JsonResponse(w, message, http.StatusOK)
}
