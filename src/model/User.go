package model

import "time"

type User struct {
	Id        *int64     `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Password  *string    `json:"password"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
