package model

import "time"

type Tag struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	User      User       `json:"user"`
	Purchases []Purchase `json:"purchases"`
	CreatedAt *time.Time `json:"createdAt"`
}
