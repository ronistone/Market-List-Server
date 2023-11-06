package model

import "time"

type Purchase struct {
	Id            *int64         `json:"id"`
	User          User           `json:"user"`
	Market        Market         `json:"market"`
	CreatedAt     *time.Time     `json:"createdAt"`
	Items         []PurchaseItem `json:"items"`
	UserId        *int64         `json:"userId"`
	MarketId      *int64         `json:"marketId"`
	TotalSpent    int64          `json:"totalSpent"`
	TotalExpected int64          `json:"totalExpected"`
}

type PurchaseItem struct {
	Id        *int64     `json:"id"`
	Purchase  *Purchase  `json:"purchase"`
	Product   Product    `json:"product"`
	Purchased bool       `json:"purchased"`
	Quantity  int        `json:"quantity"`
	Price     *int64     `json:"price"`
	CreatedAt *time.Time `json:"createdAt"`
}
