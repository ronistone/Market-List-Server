package model

import "time"

type Purchase struct {
	Id            *int64         `json:"id"`
	Name          string         `json:"name"`
	Users         []User         `json:"users"`
	Market        *Market        `json:"market"`
	CreatedAt     *time.Time     `json:"createdAt"`
	Items         []PurchaseItem `json:"items"`
	MarketId      *int64         `json:"marketId"`
	TotalSpent    int64          `json:"totalSpent"`
	TotalExpected int64          `json:"totalExpected"`
	IsFavorite    bool           `json:"isFavorite"`
	Tags          []Tag          `json:"tags"`
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
