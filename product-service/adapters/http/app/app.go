package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/build"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/logger"

	ecommerceauth "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/service/auth"

	ecommerceconfig "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/config"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/repository/products"
	product "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/services/products"
	database "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/database"
)

const (
	APP_PORT string = "APP_PORT"
)

type ProductApp struct {
	Logger         logger.Logger
	BuildInfo      model.BuildInfo
	AuthService    ecommerceauth.AuthService
	ProductService product.ProductService
	ProductInfo    string
}

type AppDependencies struct {
	logg      logger.Logger
	db        *database.DB
	jwtConfig ecommerceconfig.JWTConfig
	buildInfo model.BuildInfo
}

func initDatabase(ctx context.Context) (*database.DB, error) {
	dbConfig, err := ecommerceconfig.GetDatabaseVariables()
	if err != nil {
		return nil, err
	}

	db, err := database.InitDB(dbConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Automigrate(); err != nil {
		return nil, err
	}

	return db, nil
}

func LoadAppDependencies(ctx context.Context) (*AppDependencies, error) {
	db, err := initDatabase(ctx)
	if err != nil {
		return nil, err
	}

	jwtConfig, err := ecommerceconfig.LoadJWTConfig()
	if err != nil {
		return nil, err
	}

	buildInfo, err := build.LoadBuildInformation()
	if err != nil {
		return nil, err
	}

	logg := logger.NewLogger(logger.LogMetadata{
		"version": buildInfo.Version,
		"app":     "ecommerce/product-service",
		"from":    "product-app",
	})
	logg.SetLevelFromDefaultEnvVar()

	return &AppDependencies{
		logg:      logg,
		db:        db,
		jwtConfig: jwtConfig,
		buildInfo: buildInfo,
	}, nil
}

func NewproductApp(ctx context.Context, deps *AppDependencies) *ProductApp {

	productsRepository := products.NewProductRepository(deps.db.Conn, deps.logg)
	productservice := product.NewProductService(productsRepository, deps.logg)
	authService := ecommerceauth.NewAuthService(deps.jwtConfig.SigningKey, deps.jwtConfig.ExpiryTime, deps.jwtConfig.TokenKey, deps.logg)
	return &ProductApp{
		AuthService:    authService,
		ProductService: productservice,
		Logger:         deps.logg,
		BuildInfo:      deps.buildInfo,
		ProductInfo:    fmt.Sprintf("%v.%v:%v", deps.buildInfo.Group, deps.buildInfo.Name, deps.buildInfo.Version),
	}
}

func (app *ProductApp) Start() {
	router := app.LoadRoutes()

	port := os.Getenv(APP_PORT)
	addr := fmt.Sprintf(":%v", port)

	app.Logger.Infof("Server is running")
	app.Logger.Infof("HTTP server is listening on PORT %v", port)

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		IdleTimeout:       60 * time.Second, // idle connections
		ReadHeaderTimeout: 10 * time.Second, // request header
		ReadTimeout:       5 * time.Minute,  // request body
		WriteTimeout:      5 * time.Minute,  // response body
		MaxHeaderBytes:    1 << 20,          // 1 MB
	}

	server.ListenAndServe()

}
