package model

import "time"

type Tag struct {
	Id        *int64     `db:"id"`
	Name      string     `db:"name"`
	UserId    int64      `db:"user_id"`
	CreatedAt *time.Time `db:"created_at"`
}

type TagPurchase struct {
	PurchaseId int64 `db:"purchase_id"`
	UserId     int64 `db:"user_id"`
}
