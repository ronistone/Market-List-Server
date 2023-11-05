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
