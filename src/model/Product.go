package model

import "time"

type Product struct {
	Id        *int64     `json:"id"`
	Ean       *string    `json:"ean"`
	Name      string     `json:"name"`
	Unit      string     `json:"unit"`
	Size      int64      `json:"size"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type ProductInstance struct {
	Id        *int64     `json:"id"`
	Product   Product    `json:"product"`
	Market    *Market    `json:"market"`
	Price     int        `json:"price"`
	Precision int        `json:"precision"`
	CreatedAt *time.Time `json:"createdAt"`
}
