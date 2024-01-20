package model

import "time"

type Tag struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	UserId    int64      `json:"userid"`
	Purchases []Purchase `json:"purchases"`
	CreatedAt *time.Time `json:"createdAt"`
}
