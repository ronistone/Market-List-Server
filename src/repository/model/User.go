package model

import "time"

type User struct {
	Id        *int64     `json:"id" db:"ID"`
	Email     string     `json:"email" db:"EMAIL"`
	Name      string     `json:"name" db:"NAME"`
	Password  *string    `json:"password" db:"PASSWORD"`
	CreatedAt *time.Time `json:"createdAt" db:"CREATED_AT"`
	UpdatedAt *time.Time `json:"updatedAt" db:"UPDATED_AT"`
}
