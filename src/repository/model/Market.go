package model

import "time"

type Market struct {
	Id        *int64     `json:"id" db:"ID"`
	Name      string     `json:"name" db:"NAME"`
	CreatedAt *time.Time `json:"createdAt" db:"CREATED_AT"`
	UpdatedAt *time.Time `json:"updatedAt" db:"UPDATED_AT"`
}
