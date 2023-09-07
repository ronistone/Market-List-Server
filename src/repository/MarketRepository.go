package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/ronistone/market-list/src/model"
	"github.com/ronistone/market-list/src/util"
)

type MarketRepository interface {
	CreateMarket(ctx context.Context, market model.Market) (model.Market, error)
	UpdateMarket(ctx context.Context, market model.Market) (model.Market, error)
	DeleteMarket(ctx context.Context, id int64) error
	GetMarketById(ctx context.Context, id int64) (model.Market, error)
	List(ctx context.Context) ([]model.Market, error)
}

type Market struct {
	DbConnection *dbr.Connection
}

func CreateMarketRepository(connection *dbr.Connection) MarketRepository {
	return &Market{
		DbConnection: connection,
	}
}

func (m Market) CreateMarket(ctx context.Context, market model.Market) (model.Market, error) {
	statement := m.DbConnection.NewSession(nil).SelectBySql(`
	INSERT INTO MARKET(id, name, created_at, updated_at) 
		values (default, ?, default, default)
	RETURNING *
	`, market.Name)

	_, err := statement.LoadContext(ctx, &market)
	if err != nil {
		return model.Market{}, util.MakeErrorUnknown(err)
	}

	return market, nil
}

func (m Market) UpdateMarket(ctx context.Context, market model.Market) (model.Market, error) {
	if market.Id == nil {
		return model.Market{}, util.MakeError(util.INVALID_INPUT, "invalid Market Id")
	}
	statement := m.DbConnection.NewSession(nil).SelectBySql(`
	UPDATE MARKET SET name = ?, updated_at = NOW() 
		WHERE id = ?
	RETURNING *
	`, market.Name)

	_, err := statement.LoadContext(ctx, &market)
	if err != nil {
		return model.Market{}, util.MakeErrorUnknown(err)
	}

	return market, nil
}

func (m Market) GetMarketById(ctx context.Context, id int64) (model.Market, error) {
	statement := m.DbConnection.NewSession(nil).SelectBySql(`
	SELECT * FROM MARKET where id = ?
	`, id)

	var market model.Market
	err := statement.LoadOne(&market)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return model.Market{}, util.MakeError(util.NOT_FOUND, fmt.Sprintf("Product %s not found", id))
		}
		return model.Market{}, util.MakeErrorUnknown(err)
	}

	return market, nil
}

func (m Market) List(ctx context.Context) ([]model.Market, error) {
	statement := m.DbConnection.NewSession(nil).SelectBySql(`
	SELECT * FROM MARKET
	ORDER BY CREATED_AT DESC
-- 	LIMIT ?
	`)
	//`, limit)

	var markets []model.Market
	_, err := statement.LoadContext(ctx, &markets)
	if err != nil {
		return []model.Market{}, util.MakeErrorUnknown(err)
	}

	return markets, nil
}

func (m Market) DeleteMarket(ctx context.Context, id int64) error {
	_, err := m.DbConnection.NewSession(nil).ExecContext(ctx, `
	DELETE FROM MARKET where id = ?
	`, id)

	if err != nil {
		return util.MakeErrorUnknown(err)
	}

	return nil
}
