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

type MarketController struct {
	MarketService service.MarketService
}

func CreateMarketController(purchaseService service.MarketService) *MarketController {
	return &MarketController{
		MarketService: purchaseService,
	}
}

func (m MarketController) Register(echo *echo.Echo) error {
	v1 := echo.Group("/v1/market")
	v1.POST("/", m.CreateMarket)
	v1.PUT("/:id", m.UpdateMarket)
	v1.DELETE("/:id", m.DisableMarket)
	v1.GET("/:id", m.GetMarket)
	v1.GET("/", m.GetAllMarkets)

	return nil
}

func (m MarketController) CreateMarket(c echo.Context) error {
	var market model.Market

	if err := c.Bind(&market); err != nil {
		return handleError(c, http.StatusBadRequest, err)
	}

	market, err := m.MarketService.Create(c.Request().Context(), market)

	if err != nil {
		return handleError(c, http.StatusUnprocessableEntity, err)
	}

	return c.JSON(http.StatusCreated, market)
}

func (m MarketController) UpdateMarket(c echo.Context) error {
	var market model.Market

	if err := c.Bind(&market); err != nil {
		return handleError(c, http.StatusBadRequest, err)
	}

	id := c.Param("id")
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Market Id"))
	}
	market.Id = &idValue

	market, err = m.MarketService.Update(c.Request().Context(), market)

	if err != nil {
		return handleError(c, http.StatusUnprocessableEntity, err)
	}

	return c.JSON(http.StatusCreated, market)
}

func (m MarketController) DisableMarket(c echo.Context) error {
	id := c.Param("id")
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Market Id"))
	}

	err = m.MarketService.Disable(c.Request().Context(), idValue)

	if err != nil {
		var mkError *util.MarketListError
		if errors.As(err, &mkError) && mkError.ErrorType == util.NOT_FOUND {
			return handleError(c, http.StatusNotFound, mkError)
		}
		return handleError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

func (m MarketController) GetMarket(c echo.Context) error {
	id := c.Param("id")
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Market Id"))
	}

	products, err := m.MarketService.GetById(c.Request().Context(), idValue)

	if err != nil {
		var mkError *util.MarketListError
		if errors.As(err, &mkError) && mkError.ErrorType == util.NOT_FOUND {
			return handleError(c, http.StatusNotFound, mkError)
		}
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, products)
}

func (m MarketController) GetAllMarkets(c echo.Context) error {
	products, err := m.MarketService.List(c.Request().Context())

	if err != nil {
		var mkError *util.MarketListError
		if errors.As(err, &mkError) && mkError.ErrorType == util.NOT_FOUND {
			return handleError(c, http.StatusNotFound, mkError)
		}
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, products)
}
