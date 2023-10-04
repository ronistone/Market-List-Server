package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/ronistone/market-list/src/model"
	repositoryModel "github.com/ronistone/market-list/src/repository/model"
	"github.com/ronistone/market-list/src/util"
)

type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, purchase model.Purchase) (model.Purchase, error)
	DeletePurchase(ctx context.Context, id int64) error
	AddPurchaseItem(ctx context.Context, purchaseId int64, item model.PurchaseItem) (model.Purchase, error)
	RemovePurchaseItem(ctx context.Context, purchaseId int64, itemId int64) (model.Purchase, error)
	UpdatePurchaseItem(ctx context.Context, purchaseId int64, itemId int64, item model.PurchaseItem) error
	GetPurchaseById(ctx context.Context, id int64) (model.Purchase, error)
	GetPurchaseByIdFetchItems(ctx context.Context, id int64) (model.Purchase, error)
	GetPurchaseItemById(ctx context.Context, purchaseId int64, id int64) (model.PurchaseItem, error)
	ListPurchase(ctx context.Context, userId int64) ([]model.Purchase, error)
}

type Purchase struct {
	DbConnection *dbr.Connection
}

const (
	FETCH_PURCHASE_ITEM = `SELECT pi.id purchase_item_id,
       poi.id prod_instance_id,
       p.id prod_id,
       poi.price prod_price,
       poi.precision prod_precision,
       poi.created_at prod_inst_created_at,
       p.name prod_name,
       p.ean prod_ean,
       p.unit prod_unit,
       p.size prod_size,
       p.created_at prod_created_at,
       p.updated_at prod_updated_at,
       pi.purchased purchase_item_purchased,
       pi.quantity purchase_item_quantity,
       m.id market_id,
       m.name market_name,
       m.created_at market_created_at,
       m.updated_at market_update_at,
       m.enabled market_enabled

FROM purchase_item pi, product_instance poi, product p, market m
    where pi.product_instance_id = poi.id
	AND p.id = poi.product_id
	AND poi.market_id = m.id
`
	FETCH_PURCHASE = `SELECT
		p.id purchase_id,
		p.created_at purchase_created_at,
		mu.id market_user_id,
		mu.name market_user_name,
		-- 		    mu.email market_user_email,
		mu.created_at market_user_created_at,
		mu.updated_at market_user_updated_at,
		m.id _market_id,
		m.name market_name,
		m.created_at market_created_at,
		m.updated_at market_updated_at
	FROM purchase p, market_user mu, market m
	WHERE p.user_id = mu.id
	AND p.market_id = m.id 
`
)

func CreatePurchaseRepository(connection *dbr.Connection) PurchaseRepository {
	return &Purchase{
		DbConnection: connection,
	}
}

func (p Purchase) CreatePurchase(ctx context.Context, purchase model.Purchase) (model.Purchase, error) {
	session := p.DbConnection.NewSession(nil)
	tx, err := session.Begin()
	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	statement := tx.InsertBySql(`
	INSERT INTO PURCHASE(CREATED_AT, USER_ID, MARKET_ID) 
		values (default, ?, ?)
	RETURNING *
	`, purchase.User.Id, purchase.Market.Id)

	err = statement.LoadContext(ctx, &purchase)

	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	err = tx.Commit()
	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	return purchase, nil
}

func (p Purchase) DeletePurchase(ctx context.Context, id int64) error {
	statement := p.DbConnection.NewSession(nil).DeleteBySql(`
	DELETE FROM purchase where id = ?
	`, id)

	_, err := statement.ExecContext(ctx)
	if err != nil {
		return util.MakeErrorUnknown(err)
	}

	return nil
}

func (p Purchase) UpdatePurchaseItem(ctx context.Context, purchaseId int64, itemId int64, item model.PurchaseItem) error {
	statement := p.DbConnection.NewSession(nil).DeleteBySql(`
	UPDATE PURCHASE_ITEM SET product_instance_id = ?, purchased = ?, quantity = ? where purchase_id = ? AND id = ?
	`, item.ProductInstance.Id, item.Purchased, item.Quantity, purchaseId, itemId)

	_, err := statement.ExecContext(ctx)
	if err != nil {
		return util.MakeErrorUnknown(err)
	}

	return nil
}

func (p Purchase) AddPurchaseItem(ctx context.Context, purchaseId int64, item model.PurchaseItem) (model.Purchase, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	INSERT INTO PURCHASE_ITEM(ID, PRODUCT_INSTANCE_ID, PURCHASE_ID, QUANTITY) 
	values (default, ?, ?, ?)`, *item.ProductInstance.Id, purchaseId, item.Quantity)

	_, err := statement.LoadContext(ctx, &item)
	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	return p.GetPurchaseById(ctx, purchaseId)
}

func (p Purchase) RemovePurchaseItem(ctx context.Context, purchaseId int64, itemId int64) (model.Purchase, error) {
	statement := p.DbConnection.NewSession(nil).DeleteBySql(`
	DELETE FROM purchase_item where id = ? AND purchase_id = ?
	`, itemId, purchaseId)

	_, err := statement.ExecContext(ctx)

	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	return p.GetPurchaseByIdFetchItems(ctx, purchaseId)
}

func (p Purchase) getAllPurchaseItemByPurchaseId(ctx context.Context, purchaseId int64, purchase model.Purchase) ([]model.PurchaseItem, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(FETCH_PURCHASE_ITEM+`
	AND pi.purchase_id = ?
	ORDER BY purchase_item_purchased, purchase_item_quantity ASC
	`, purchaseId)

	var items []repositoryModel.PurchaseItemProductInstance
	_, err := statement.LoadContext(ctx, &items)
	if err != nil {
		return []model.PurchaseItem{}, util.MakeErrorUnknown(err)
	}

	var results = make([]model.PurchaseItem, len(items))

	for i, v := range items {
		results[i] = v.ToPurchaseItem()
		results[i].Purchase = &purchase
		results[i].ProductInstance.Market = &purchase.Market
	}

	return results, nil
}

func (p Purchase) GetPurchaseByIdFetchItems(ctx context.Context, id int64) (model.Purchase, error) {
	return p.getPurchaseByIdInternal(ctx, id, true)
}

func (p Purchase) GetPurchaseById(ctx context.Context, id int64) (model.Purchase, error) {
	return p.getPurchaseByIdInternal(ctx, id, false)
}

func (p Purchase) getPurchaseByIdInternal(ctx context.Context, id int64, fetchItems bool) (model.Purchase, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(FETCH_PURCHASE+`
		  AND p.ID = ?
		`, id)

	var purchase repositoryModel.PurchaseEntity
	err := statement.LoadOne(&purchase)

	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return model.Purchase{}, util.MakeError(util.NOT_FOUND, fmt.Sprintf("Purchase %d not found", id))
		}
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	result := purchase.ToPurchase()
	if !fetchItems {
		return result, nil
	}
	result.Items, err = p.getAllPurchaseItemByPurchaseId(ctx, id, result)
	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	return result, nil
}

func (p Purchase) GetPurchaseItemById(ctx context.Context, purchaseId int64, id int64) (model.PurchaseItem, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(FETCH_PURCHASE_ITEM+`
	AND pi.purchase_id = ?
	AND pi.id = ?
	ORDER BY purchase_item_purchased, purchase_item_quantity ASC
	`, purchaseId, id)

	var item repositoryModel.PurchaseItemProductInstance
	err := statement.LoadOne(&item)
	if err != nil {
		return model.PurchaseItem{}, util.MakeErrorUnknown(err)
	}

	return item.ToPurchaseItem(), nil
}

func (p Purchase) ListPurchase(ctx context.Context, userId int64) ([]model.Purchase, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(FETCH_PURCHASE+`
	AND USER_ID = ?
	`, userId)

	var items []repositoryModel.PurchaseEntity
	_, err := statement.LoadContext(ctx, &items)
	if err != nil {
		return []model.Purchase{}, util.MakeErrorUnknown(err)
	}

	results := make([]model.Purchase, len(items))

	for i, v := range items {
		results[i] = v.ToPurchase()
	}

	return results, nil
}
