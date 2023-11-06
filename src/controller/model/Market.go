package model

import "github.com/ronistone/market-list/src/model"

type Market struct {
	Id   *int64 `json:"id"`
	Name string `json:"name"`
}

func (m *Market) FromModel(marketModel model.Market) {
	m.Id = marketModel.Id
	m.Name = marketModel.Name
}
