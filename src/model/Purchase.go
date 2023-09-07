package model

import "time"

type Purchase struct {
	Id        *int64         `json:"id"`
	User      User           `json:"user"`
	Market    Market         `json:"market"`
	CreatedAt *time.Time     `json:"createdAt"`
	Items     []PurchaseItem `json:"items"`
	UserId    *int64         `json:"userId"`
	MarketId  *int64         `json:"marketId"`
}

type PurchaseItem struct {
	Id              *int64          `json:"id"`
	Purchase        *Purchase       `json:"purchase"`
	ProductInstance ProductInstance `json:"productInstance"`
}
