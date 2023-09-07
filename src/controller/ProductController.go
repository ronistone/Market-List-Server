package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/ronistone/market-list/src/model"
	"github.com/ronistone/market-list/src/service"
	"github.com/ronistone/market-list/src/util"
	"net/http"
	"strconv"
)

type ProductController struct {
	productService service.ProductService
}

func CreateProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (p ProductController) Register(echo *echo.Echo) error {
	v1 := echo.Group("/v1/product")
	v1.GET("/ean/:ean", p.GetProductByEan)
	v1.GET("/:id", p.GetProductById)
	v1.GET("/name/:name", p.GetProductByName)
	v1.POST("/", p.CreateProduct)
	v1.PUT("/:id", p.UpdateProduct)

	return nil
}

func handleError(echo echo.Context, statusCode int, err error) error {
	//panic(err)
	return echo.JSON(statusCode, err)
}

func (p ProductController) UpdateProduct(c echo.Context) error {
	var product model.Product

	if err := c.Bind(&product); err != nil {
		return handleError(c, http.StatusBadRequest, err)
	}

	id := c.Param("id")
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Product Id"))
	}
	product.Id = &idValue

	product, err = p.productService.Update(c.Request().Context(), product)

	if err != nil {
		return handleError(c, http.StatusUnprocessableEntity, err)
	}

	return c.JSON(http.StatusCreated, product)
}

func (p ProductController) GetProductByName(c echo.Context) error {
	name := c.Param("name")

	products, err := p.productService.GetByName(c.Request().Context(), name)

	if err != nil {
		return handleError(c, http.StatusUnprocessableEntity, err)
	}

	return c.JSON(http.StatusOK, products)
}

func (p ProductController) CreateProduct(c echo.Context) error {
	var product model.Product

	if err := c.Bind(&product); err != nil {
		return handleError(c, http.StatusBadRequest, err)
	}

	product, err := p.productService.Create(c.Request().Context(), product)

	if err != nil {
		return handleError(c, http.StatusUnprocessableEntity, err)
	}

	return c.JSON(http.StatusCreated, product)
}

func (p ProductController) GetProductById(c echo.Context) error {
	id := c.Param("id")
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Product Id"))
	}

	products, err := p.productService.GetById(c.Request().Context(), idValue)

	if err != nil {
		var mkError *util.MarketListError
		if errors.As(err, &mkError) && mkError.ErrorType == util.NOT_FOUND {
			return handleError(c, http.StatusNotFound, mkError)
		}
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, products)
}

func (p ProductController) GetProductByEan(c echo.Context) error {
	ean := c.Param("ean")

	products, err := p.productService.GetByEan(c.Request().Context(), ean)

	if err != nil {
		return handleError(c, http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, products)
}
