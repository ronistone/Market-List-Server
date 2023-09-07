package main

import (
	"context"
	"errors"
	"github.com/gocraft/dbr/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/ronistone/market-list/src/controller"
	"github.com/ronistone/market-list/src/repository"
	"github.com/ronistone/market-list/src/service"
	"github.com/ronistone/market-list/src/util"
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
		if err := e.Start("0.0.0.0:8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("Shutting down the server!")
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
	e := ConfigureServer()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("logger", &util.ContextLogger{Logger: c.Logger()})
			return next(c)
		}
	})

	db, err := dbr.Open("postgres", "host=localhost port=5432 user=postgres password='market-list' dbname=market_list sslmode=disable timezone=UTC", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(20)

	userRepository := repository.CreateUserRepository(db)
	userService := service.CreateUserService(userRepository)

	productRepository := repository.CreateProductRepository(db)
	productInstanceRepository := repository.CreateProductInstanceRepository(db)
	productService := service.CreateProductService(productRepository, productInstanceRepository)
	productController := controller.CreateProductController(productService)

	marketRepository := repository.CreateMarketRepository(db)
	marketService := service.CreateMarketService(marketRepository)
	marketController := controller.CreateMarketController(marketService)

	purchaseRepository := repository.CreatePurchaseRepository(db)
	purchaseService := service.CreatePurchaseService(purchaseRepository, productService, marketService, userService)
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
