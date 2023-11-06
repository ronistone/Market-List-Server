package model

import (
	"github.com/ronistone/market-list/src/model"
)

type Product struct {
	Id   *int64  `json:"id"`
	Ean  *string `json:"ean"`
	Name string  `json:"name"`
	Unit string  `json:"unit"`
	Size int64   `json:"size"`
}

func (p *Product) FromModel(productModel model.Product) {
	p.Id = productModel.Id
	p.Ean = productModel.Ean
	p.Name = productModel.Name
	p.Unit = productModel.Unit
	p.Size = productModel.Size
}
