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

type PurchaseController struct {
	PurchaseService service.PurchaseService
}

func CreatePurchaseController(purchaseService service.PurchaseService) *PurchaseController {
	return &PurchaseController{
		PurchaseService: purchaseService,
	}
}

func (p PurchaseController) Register(echo *echo.Echo) error {
	v1 := echo.Group("/v1/purchase")
	v1.POST("/", p.CreatePurchase)
	v1.DELETE("/:id", p.DeletePurchase)
	v1.POST("/:id/item/", p.AddItem)
	v1.DELETE("/:id/item/:itemId", p.RemoveItem)
	v1.GET("/:id", p.GetPurchase)
	v1.GET("/", p.GetAllPurchase)

	return nil
}

func (p PurchaseController) CreatePurchase(c echo.Context) error {
	var purchase model.Purchase

	if err := c.Bind(&purchase); err != nil {
		return handleError(c, http.StatusBadRequest, err)
	}

	purchase, err := p.PurchaseService.CreatePurchase(c.Request().Context(), purchase)

	if err != nil {
		return handleError(c, http.StatusUnprocessableEntity, err)
	}

	return c.JSON(http.StatusCreated, purchase)
}

func (p PurchaseController) AddItem(c echo.Context) error {
	id := c.Param("id")
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Purchase Id"))
	}

	var purchaseItem model.PurchaseItem
	if err := c.Bind(&purchaseItem); err != nil {
		return handleError(c, http.StatusBadRequest, err)
	}

	purchase, err := p.PurchaseService.AddItem(c.Request().Context(), idValue, purchaseItem)

	if err != nil {
		var mkError *util.MarketListError
		if errors.As(err, &mkError) && mkError.ErrorType == util.NOT_FOUND {
			return handleError(c, http.StatusNotFound, mkError)
		}
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, purchase)
}

func (p PurchaseController) RemoveItem(c echo.Context) error {
	id := c.Param("id")
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Purchase Id"))
	}

	itemId := c.Param("itemId")
	itemIdValue, err := strconv.ParseInt(itemId, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Purchase Item Id"))
	}

	purchase, err := p.PurchaseService.RemoveItem(c.Request().Context(), idValue, itemIdValue)

	if err != nil {
		var mkError *util.MarketListError
		if errors.As(err, &mkError) && mkError.ErrorType == util.NOT_FOUND {
			return handleError(c, http.StatusNotFound, mkError)
		}
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, purchase)
}

func (p PurchaseController) GetPurchase(c echo.Context) error {
	id := c.Param("id")
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Purchase Id"))
	}

	products, err := p.PurchaseService.GetPurchase(c.Request().Context(), idValue)

	if err != nil {
		var mkError *util.MarketListError
		if errors.As(err, &mkError) && mkError.ErrorType == util.NOT_FOUND {
			return handleError(c, http.StatusNotFound, mkError)
		}
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, products)
}

func (p PurchaseController) GetAllPurchase(c echo.Context) error {

	purchases, err := p.PurchaseService.GetAllPurchase(c.Request().Context())

	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, purchases)
}

func (p PurchaseController) DeletePurchase(c echo.Context) error {
	id := c.Param("id")
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return handleError(c, http.StatusBadRequest, util.MakeError(util.INVALID_INPUT, "invalid Purchase Id"))
	}

	err = p.PurchaseService.DeletePurchase(c.Request().Context(), idValue)

	if err != nil {
		var mkError *util.MarketListError
		if errors.As(err, &mkError) && mkError.ErrorType == util.NOT_FOUND {
			return handleError(c, http.StatusNotFound, mkError)
		}
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, nil)
}
