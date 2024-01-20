package main

import (
	"context"
	"errors"
	"github.com/gocraft/dbr/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/ronistone/market-list/src/config"
	"github.com/ronistone/market-list/src/controller"
	myMiddleware "github.com/ronistone/market-list/src/middleware"
	"github.com/ronistone/market-list/src/repository"
	"github.com/ronistone/market-list/src/service"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func ConfigureServer() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
	return e
}

func GracefullyStart(e *echo.Echo) {
	go func() {
		tlsEnabled := config.GetTlsEnabled()
		var err error

		if tlsEnabled {
			crtPath := config.GetTlsCrtPath()
			keyPath := config.GetTlsKeyPath()
			err = e.StartTLS("0.0.0.0:8080", crtPath, keyPath)
		} else {
			err = e.Start("0.0.0.0:8080")
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("Shutting down the server! %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	} else {
		e.Logger.Fatal("Graceful Shutting down the server!")
	}
}

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	e := ConfigureServer()

	e.Use(myMiddleware.InjectLogger)
	e.Use(myMiddleware.InjectUserId)

	db, err := dbr.Open("postgres", config.GetDatabaseDSN(), nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(20)

	userRepository := repository.CreateUserRepository(db)
	userService := service.CreateUserService(userRepository)

	productRepository := repository.CreateProductRepository(db)
	productService := service.CreateProductService(productRepository)
	productController := controller.CreateProductController(productService)

	marketRepository := repository.CreateMarketRepository(db)
	marketService := service.CreateMarketService(marketRepository)
	marketController := controller.CreateMarketController(marketService)

	purchaseRepository := repository.CreatePurchaseRepository(db)
	purchaseService := service.CreatePurchaseService(purchaseRepository, productService, userService)
	purchaseController := controller.CreatePurchaseController(purchaseService)

	err = productController.Register(e)
	if err != nil {
		panic(err)
	}

	err = marketController.Register(e)
	if err != nil {
		panic(err)
	}

	err = purchaseController.Register(e)
	if err != nil {
		panic(err)
	}

	GracefullyStart(e)
}
