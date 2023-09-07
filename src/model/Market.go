package model

import "time"

type Market struct {
	Id        *int64     `json:"id"`
	Name      string     `json:"name"`
	Enabled   bool       `json:"enabled"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}
