package model

import "time"

type Product struct {
	Id        *int64     `json:"id" db:"ID"`
	Ean       *string    `json:"ean" db:"EAN"`
	Name      string     `json:"name" db:"NAME"`
	Unit      string     `json:"unit" db:"UNIT"`
	Size      int64      `json:"size" db:"SIZE"`
	CreatedAt *time.Time `json:"createdAt" db:"CREATED_AT"`
	UpdatedAt *time.Time `json:"updatedAt" db:"UPDATED_AT"`
}

type ProductInstance struct {
	Id        *int64     `json:"id" db:"ID"`
	ProductId int64      `json:"productId" db:"PRODUCT_ID"`
	MarketId  int64      `json:"marketId" db:"MARKET_ID"`
	Price     *int       `json:"price" db:"PRICE"`
	Precision int        `json:"precision" db:"PRECISION"`
	CreatedAt *time.Time `json:"createdAt" db:"CREATED_AT"`
}
