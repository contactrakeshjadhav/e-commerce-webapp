package app

import (
	"net/http"

	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/middleware"

	"github.com/bsm/openmetrics"
	"github.com/bsm/openmetrics/omhttp"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/adapters/http/controllers"
	authmiddleware "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/adapters/http/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	openapi "github.com/go-openapi/runtime/middleware"
)

const (
	//general routes
	ROUTE_INDEX        = "/"
	ROUTE_PING         = "/ping"
	ROUTE_BUILD_INFO   = "/getBuildInfo"
	ROUTE_SWAGGER_UI   = "/swagger-ui"
	ROUTE_SWAGGER_FILE = "/swagger.json"

	// product routes
	ROUTE_SAVE_PRODUCT                 = "/v1/saveProduct"
	ROUTE_DELETE_PRODUCT               = "/v1/deleteProduct"
	ROUTE_GET_PRODUCTS                 = "/v1/getProducts"
	ROUTE_GET_PRODUCTS_BY_PARTIAL_NAME = "/v1/getProductsByPartialName"
	ROUTE_GET_PRODUCTS_BY_IDS          = "/v1/getProductsByIDs"
	ROUTE_GET_PRODUCTS_BY_OBJECT_IDS   = "/v1/getProductsByObjectIDs"

	// object products routes
	ROUTE_ADD_PRODUCTS_TO_OBJECT      = "/v1/addProductsToObject"
	ROUTE_REMOVE_PRODUCTS_FROM_OBJECT = "/v1/removeProductsFromObject"
	ROUTE_GET_PRODUCT                 = "/v1/getProduct"
	ROUTE_FIND_OBJECTS_BY_PRODUCTS    = "/v1/findObjectItemsByProducts"
	ROUTE_UPDATE_PRODUCTS_TO_OBJECT   = "/v1/updateProductsToObject"
)

// LoadRoutes this is responsible to create a chi
// router and load all the application's routes
// return a *chi.Mux
func (app *ProductApp) LoadRoutes() *chi.Mux {
	// init router
	router := chi.NewRouter()

	// common middlewares
	router.Use(middleware.Recoverer)

	// setup required controllers
	baseController := controllers.NewBaseController(app.BuildInfo, app.Logger)
	productController := controllers.NewProductController(app.ProductService, app.Logger)

	swaggerOpts := openapi.SwaggerUIOpts{SpecURL: "/product-api/swagger.json", BasePath: "/", Path: "swagger-ui"}
	swaggerUI := openapi.SwaggerUI(swaggerOpts, nil)

	authMiddleware := authmiddleware.NewAuthorizationMiddleware(app.Logger)

	// setup open routes
	router.Get(ROUTE_INDEX, baseController.Index)
	router.Get(ROUTE_PING, baseController.Ping)
	router.Post(ROUTE_BUILD_INFO, baseController.GetBuildInfo)

	// setup swagger routes
	router.Handle(ROUTE_SWAGGER_UI, swaggerUI)
	router.HandleFunc(ROUTE_SWAGGER_FILE, controllers.SwaggerDocsHandler)

	// protected routes
	router.Group(func(r chi.Router) {
		// authorization middleware
		r.Use(app.AuthService.ContextInitiator)
		r.Use(app.AuthService.Authenticator(app.AuthService.GetTokenFromHeader))

		//The following code block adds /metrics endpoint middleware
		//	this middleware tracks HTTP requests and response times
		//	it also info logs the results of each request
		metricsRegistry := openmetrics.DefaultRegistry()
		requestCount := metricsRegistry.Counter(openmetrics.Desc{
			Name:   "http_request",
			Help:   "http request count",
			Labels: []string{"status", "correlationID", "user", "service", "method", "path"},
		})
		responseTime := metricsRegistry.Histogram(openmetrics.Desc{
			Name:   "http_request",
			Unit:   "seconds",
			Help:   "http response time histogram",
			Labels: []string{"status", "correlationID", "user", "service", "method", "path"},
		}, []float64{.005, .01, .05, .1, .5, 1, 5, 10})

		r.Use(commons.Metrics(requestCount, responseTime, app.Logger, app.ProductInfo))
		r.Handle("/metrics", omhttp.NewHandler(metricsRegistry))
		// End of /metrics endpoint code

		// products
		r.Post(ROUTE_GET_PRODUCTS, productController.GetProducts)
		r.Post(ROUTE_GET_PRODUCTS_BY_PARTIAL_NAME, productController.GetProductsByPartialName)
		r.Post(ROUTE_GET_PRODUCT, productController.GetProduct)
		r.Post(ROUTE_GET_PRODUCTS_BY_IDS, productController.GetProductsByIDs)
		r.Post(ROUTE_SAVE_PRODUCT, authMiddleware.ValidateAdminAccess(
			http.HandlerFunc(productController.SaveProduct),
		))
		r.Post(ROUTE_DELETE_PRODUCT, authMiddleware.ValidateAdminAccess(
			http.HandlerFunc(productController.DeleteProduct),
		))

	})
	return router
}
