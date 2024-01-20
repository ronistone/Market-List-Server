package model

import (
	"github.com/ronistone/market-list/src/model"
	"time"
)

type Purchase struct {
	Id            *int64         `json:"id,omitempty"`
	Name          string         `json:"name"`
	User          []User         `json:"user,omitempty"`
	Market        *Market        `json:"market,omitempty"`
	CreatedAt     *time.Time     `json:"createdAt,omitempty"`
	Items         []PurchaseItem `json:"items,omitempty"`
	TotalSpent    int64          `json:"totalSpent"`
	TotalExpected int64          `json:"totalExpected"`
	IsFavorite    bool           `json:"isFavorite"`
	Tags          []Tag          `json:"tags,omitempty"`
}

type PurchaseItem struct {
	Id        *int64     `json:"id"`
	Product   Product    `json:"product"`
	Purchased bool       `json:"purchased"`
	Quantity  int        `json:"quantity"`
	Price     *int64     `json:"price"`
	CreatedAt *time.Time `json:"createdAt"`
}

func (pi *PurchaseItem) FromModel(itemModel model.PurchaseItem) {
	pi.Id = itemModel.Id

	product := Product{}
	product.FromModel(itemModel.Product)
	pi.Product = product
	pi.Purchased = itemModel.Purchased
	pi.Quantity = itemModel.Quantity
	pi.Price = itemModel.Price
	pi.CreatedAt = itemModel.CreatedAt
}

func (p *Purchase) FromModel(purchaseModel model.Purchase) {

	p.Id = purchaseModel.Id

	var users []User

	for _, user := range purchaseModel.Users {
		users = append(users, User{
			Id:    *user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	var tags []Tag
	for _, tag := range purchaseModel.Tags {
		tags = append(tags, Tag{Id: tag.Id, Name: tag.Name})
	}

	p.User = users

	if purchaseModel.Market != nil {
		market := Market{}
		market.FromModel(*purchaseModel.Market)
		p.Market = &market
	}

	p.CreatedAt = purchaseModel.CreatedAt

	if purchaseModel.Items != nil {
		items := make([]PurchaseItem, len(purchaseModel.Items), len(purchaseModel.Items))
		for i, _ := range purchaseModel.Items {
			items[i].FromModel(purchaseModel.Items[i])
		}
		p.Items = items
	}
	p.TotalSpent = purchaseModel.TotalSpent
	p.TotalExpected = purchaseModel.TotalExpected
	p.Name = purchaseModel.Name
	p.IsFavorite = purchaseModel.IsFavorite

}
