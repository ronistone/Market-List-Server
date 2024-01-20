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
	DeletePurchase(ctx context.Context, userId, id int64) error
	AddPurchaseItem(ctx context.Context, userId, purchaseId int64, item model.PurchaseItem) (model.Purchase, error)
	RemovePurchaseItem(ctx context.Context, userId, purchaseId int64, itemId int64) (model.Purchase, error)
	UpdatePurchaseItem(ctx context.Context, userId, purchaseId, itemId int64, item model.PurchaseItem) error
	GetPurchaseById(ctx context.Context, userId, id int64) (model.Purchase, error)
	GetPurchaseByIdFetchItems(ctx context.Context, userId, id int64) (model.Purchase, error)
	GetPurchaseItemById(ctx context.Context, userId, purchaseId int64, id int64) (model.PurchaseItem, error)
	ListPurchase(ctx context.Context, userId int64) ([]model.Purchase, error)
}

type Purchase struct {
	DbConnection *dbr.Connection
}

const (
	FETCH_PURCHASE_ITEM = `SELECT pi.id purchase_item_id,
       pi.purchased purchase_item_purchased,
       pi.quantity purchase_item_quantity,
       pi.created_at purchase_item_created_at,
       pi.price purchase_item_price,
       p.id prod_id,
       p.name prod_name,
       p.ean prod_ean,
       p.unit prod_unit,
       p.size prod_size,
       p.created_at prod_created_at,
       p.updated_at prod_updated_at
FROM purchase_item pi 
    INNER JOIN product p ON p.id = pi.product_id
	INNER JOIN purchase_user pu ON pu.purchase_id = pi.purchase_id AND pu.user_id = ?
    where 1=1
`
	FETCH_PURCHASE = `SELECT
		p.id purchase_id,
		p.created_at purchase_created_at,
		p.name purchase_name,
		p.is_favorite purchase_is_favorite,
		m.id _market_id,
		m.name market_name,
		m.created_at market_created_at,
		m.updated_at market_updated_at
	FROM purchase p
	    INNER JOIN purchase_user pu on pu.user_id = ? AND pu.purchase_id = p.id
	    LEFT JOIN market m ON p.market_id = m.id
	WHERE 1=1
`
)

func CreatePurchaseRepository(connection *dbr.Connection) PurchaseRepository {
	return &Purchase{
		DbConnection: connection,
	}
}

func relateUsersToPurchase(ctx context.Context, tx *dbr.Tx, purchaseId int64, users []model.User) error {
	for _, user := range users {
		statement := tx.InsertBySql(`
			INSERT INTO PURCHASE_USER (purchase_id, user_id) VALUES (?, ?)
		`, purchaseId, user.Id)
		_, err := statement.ExecContext(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Purchase) CreatePurchase(ctx context.Context, purchase model.Purchase) (model.Purchase, error) {
	session := p.DbConnection.NewSession(nil)
	tx, err := session.Begin()
	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	statement := tx.InsertBySql(`
	INSERT INTO PURCHASE(CREATED_AT, NAME, MARKET_ID, IS_FAVORITE) 
		values (default, ?, ?, ?)
	RETURNING *
	`, purchase.Name, purchase.MarketId, purchase.IsFavorite)

	err = statement.LoadContext(ctx, &purchase)

	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	err = relateUsersToPurchase(ctx, tx, *purchase.Id, purchase.Users)

	if err != nil {
		_ = tx.Rollback()
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	err = tx.Commit()
	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	return purchase, nil
}

func (p Purchase) DeletePurchase(ctx context.Context, userId, id int64) error {
	statement := p.DbConnection.NewSession(nil).DeleteBySql(`
	DELETE FROM purchase where id = (
		SELECT pu.purchase_id FROM purchase_user pu WHERE pu.purchase_id = ? AND pu.user_id = ?
	)
	`, id, userId)

	_, err := statement.ExecContext(ctx)
	if err != nil {
		return util.MakeErrorUnknown(err)
	}

	return nil
}

func (p Purchase) UpdatePurchaseItem(ctx context.Context, userId, purchaseId int64, itemId int64, item model.PurchaseItem) error {
	statement := p.DbConnection.NewSession(nil).DeleteBySql(`
	UPDATE PURCHASE_ITEM pi
		SET purchased = ?, quantity = ?, price = ?, product_id = ?
		FROM purchase_user pu
		WHERE pi.id = ? AND pi.purchase_id = pu.purchase_id AND pu.user_id = ? AND pi.purchase_id = ?
	`, item.Purchased, item.Quantity, item.Price, item.Product.Id, itemId, userId, purchaseId)

	_, err := statement.ExecContext(ctx)
	if err != nil {
		return util.MakeErrorUnknown(err)
	}

	return nil
}

func (p Purchase) AddPurchaseItem(ctx context.Context, userId, purchaseId int64, item model.PurchaseItem) (model.Purchase, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	INSERT INTO PURCHASE_ITEM(ID, PURCHASE_ID, PRODUCT_ID, QUANTITY, PRICE) 
	values (default, ?, ?, ?, ?)`, purchaseId, item.Product.Id, item.Quantity, item.Price)

	_, err := statement.LoadContext(ctx, &item)
	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	return p.GetPurchaseById(ctx, userId, purchaseId)
}

func (p Purchase) RemovePurchaseItem(ctx context.Context, userId, purchaseId int64, itemId int64) (model.Purchase, error) {
	statement := p.DbConnection.NewSession(nil).DeleteBySql(`
	DELETE FROM purchase_item pi
	USING purchase_user pu
	WHERE pi.id = ? AND pi.purchase_id = pu.purchase_id AND pu.user_id = ? AND pi.purchase_id = ?
	`, itemId, userId, purchaseId)

	_, err := statement.ExecContext(ctx)

	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	return p.GetPurchaseByIdFetchItems(ctx, userId, purchaseId)
}

func (p Purchase) getAllPurchaseItemByPurchaseId(ctx context.Context, userId, purchaseId int64, purchase model.Purchase) ([]model.PurchaseItem, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(FETCH_PURCHASE_ITEM+`
	AND pi.purchase_id = ?
	ORDER BY purchase_item_purchased, purchase_item_quantity ASC
	`, userId, purchaseId)

	var items []repositoryModel.PurchaseItemProductInstance
	_, err := statement.LoadContext(ctx, &items)
	if err != nil {
		return []model.PurchaseItem{}, util.MakeErrorUnknown(err)
	}

	var results = make([]model.PurchaseItem, len(items))

	for i, v := range items {
		results[i] = v.ToPurchaseItem()
		results[i].Purchase = &purchase
	}

	return results, nil
}

func (p Purchase) GetPurchaseByIdFetchItems(ctx context.Context, userId, id int64) (model.Purchase, error) {
	return p.getPurchaseByIdInternal(ctx, userId, id, true)
}

func (p Purchase) GetPurchaseById(ctx context.Context, userId, id int64) (model.Purchase, error) {
	return p.getPurchaseByIdInternal(ctx, userId, id, false)
}

func (p Purchase) getPurchaseByIdInternal(ctx context.Context, userId, id int64, fetchItems bool) (model.Purchase, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(FETCH_PURCHASE+`
		  AND p.ID = ?
		`, userId, id)

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
	result.Items, err = p.getAllPurchaseItemByPurchaseId(ctx, userId, id, result)
	if err != nil {
		return model.Purchase{}, util.MakeErrorUnknown(err)
	}

	return result, nil
}

func (p Purchase) GetPurchaseItemById(ctx context.Context, userId, purchaseId int64, id int64) (model.PurchaseItem, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(FETCH_PURCHASE_ITEM+`
	AND pi.purchase_id = ?
	AND pi.id = ?
	ORDER BY purchase_item_purchased, purchase_item_quantity ASC
	`, userId, purchaseId, id)

	var item repositoryModel.PurchaseItemProductInstance
	err := statement.LoadOne(&item)
	if err != nil {
		return model.PurchaseItem{}, util.MakeErrorUnknown(err)
	}

	return item.ToPurchaseItem(), nil
}

func (p Purchase) ListPurchase(ctx context.Context, userId int64) ([]model.Purchase, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(FETCH_PURCHASE, userId)

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
