package model

import (
	"github.com/ronistone/market-list/src/model"
	"time"
)

type PurchaseEntity struct {
	Id              *int64     `db:"purchase_id"`
	CreatedAt       *time.Time `db:"purchase_created_at"`
	UserId          *int64     `db:"market_user_id"`
	UserName        string     `db:"market_user_name"`
	Email           *string    `db:"EMAIL"`
	UserCreatedAt   *time.Time `db:"market_user_created_at"`
	UserUpdatedAt   *time.Time `db:"market_user_updated_at"`
	MarketId        *int64     `db:"_market_id"`
	MarketName      string     `db:"market_name"`
	MarketCreatedAt *time.Time `db:"market_created_at"`
	MarketUpdatedAt *time.Time `db:"market_updated_at"`
}

func (p PurchaseEntity) ToPurchase() model.Purchase {
	return model.Purchase{
		Id: p.Id,
		User: model.User{
			Id:        p.UserId,
			Email:     "",
			Name:      p.UserName,
			CreatedAt: p.UserCreatedAt,
			UpdatedAt: p.UserUpdatedAt,
		},
		Market: model.Market{
			Id:        p.MarketId,
			Name:      p.MarketName,
			CreatedAt: p.MarketCreatedAt,
			UpdatedAt: p.MarketUpdatedAt,
		},
		CreatedAt: p.CreatedAt,
		Items:     nil,
		UserId:    p.UserId,
		MarketId:  p.MarketId,
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
