package service

import (
	"context"
	"github.com/ronistone/market-list/src/model"
	"github.com/ronistone/market-list/src/repository"
)

type MarketService interface {
	Create(ctx context.Context, market model.Market) (model.Market, error)
	Update(ctx context.Context, market model.Market) (model.Market, error)
	GetById(ctx context.Context, id int64) (model.Market, error)
	List(ctx context.Context) ([]model.Market, error)
	Disable(ctx context.Context, id int64) error
}

type Market struct {
	MarketRepository repository.MarketRepository
}

func CreateMarketService(marketRepository repository.MarketRepository) MarketService {
	return &Market{
		MarketRepository: marketRepository,
	}
}

func (m Market) Create(ctx context.Context, market model.Market) (model.Market, error) {
	return m.MarketRepository.CreateMarket(ctx, market)
}

func (m Market) Update(ctx context.Context, market model.Market) (model.Market, error) {
	return m.MarketRepository.UpdateMarket(ctx, market)
}

func (m Market) GetById(ctx context.Context, id int64) (model.Market, error) {
	return m.MarketRepository.GetMarketById(ctx, id)
}

func (m Market) List(ctx context.Context) ([]model.Market, error) {
	return m.MarketRepository.List(ctx)
}

func (m Market) Disable(ctx context.Context, id int64) error {
	market, err := m.MarketRepository.GetMarketById(ctx, id)
	if err != nil {
		return err
	}

	market.Enabled = false

	_, err = m.MarketRepository.UpdateMarket(ctx, market)
	return err
}
