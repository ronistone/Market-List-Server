package model

import (
	"github.com/ronistone/market-list/src/model"
	"time"
)

type PurchaseEntity struct {
	Id              *int64     `db:"purchase_id"`
	Name            string     `db:"purchase_name"`
	IsFavorite      bool       `db:"purchase_is_favorite"`
	CreatedAt       *time.Time `db:"purchase_created_at"`
	MarketId        *int64     `db:"_market_id"`
	MarketName      *string    `db:"market_name"`
	MarketCreatedAt *time.Time `db:"market_created_at"`
	MarketUpdatedAt *time.Time `db:"market_updated_at"`
}

func (p PurchaseEntity) ToPurchase() model.Purchase {
	var marketResult *model.Market
	if p.MarketId != nil && p.MarketName != nil {
		var market model.Market
		market.Id = p.MarketId
		market.Name = *p.MarketName
		market.CreatedAt = p.MarketCreatedAt
		market.UpdatedAt = p.MarketUpdatedAt
		marketResult = &market
	}

	return model.Purchase{
		Id:         p.Id,
		Name:       p.Name,
		Market:     marketResult,
		CreatedAt:  p.CreatedAt,
		Items:      nil,
		IsFavorite: p.IsFavorite,
		MarketId:   p.MarketId,
	}
}

type PurchaseItemProductInstance struct {
	PurchaseItemId        *int64     `db:"purchase_item_id"`
	PurchaseItemPurchased bool       `db:"purchase_item_purchased"`
	PurchaseItemQuantity  int        `db:"purchase_item_quantity"`
	PurchaseItemCreatedAt *time.Time `db:"purchase_item_created_at"`
	Price                 *int64     `db:"purchase_item_price"`
	ProductId             *int64     `db:"prod_id"`
	ProductName           string     `db:"prod_name"`
	ProductEan            *string    `db:"prod_ean"`
	ProductUnit           string     `db:"prod_unit"`
	ProductSize           int64      `db:"prod_size"`
	ProductCreatedAt      *time.Time `db:"prod_created_at"`
	ProductUpdatedAt      *time.Time `db:"prod_updated_at"`
}

func (p PurchaseItemProductInstance) ToPurchaseItem() model.PurchaseItem {
	return model.PurchaseItem{
		Id:       p.PurchaseItemId,
		Purchase: nil,
		Product: model.Product{
			Id:        p.ProductId,
			Ean:       p.ProductEan,
			Name:      p.ProductName,
			Unit:      p.ProductUnit,
			Size:      p.ProductSize,
			CreatedAt: p.ProductCreatedAt,
			UpdatedAt: p.ProductUpdatedAt,
		},
		Price:     p.Price,
		CreatedAt: p.PurchaseItemCreatedAt,
		Purchased: p.PurchaseItemPurchased,
		Quantity:  p.PurchaseItemQuantity,
	}
}
