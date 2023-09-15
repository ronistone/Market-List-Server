package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/ronistone/market-list/src/model"
	"github.com/ronistone/market-list/src/util"
)

type ProductInstanceRepository interface {
	CreateProduct(ctx context.Context, product model.ProductInstance) (model.ProductInstance, error)
	GetLastProductInstanceByProductId(ctx context.Context, id int64) (model.ProductInstance, error)
	GetProductInstanceByPurchase(ctx context.Context, purchaseId int64) ([]model.ProductInstance, error)
}

type ProductInstance struct {
	DbConnection *dbr.Connection
}

func CreateProductInstanceRepository(connection *dbr.Connection) ProductInstanceRepository {
	return &ProductInstance{
		DbConnection: connection,
	}
}

func (p ProductInstance) CreateProduct(ctx context.Context, product model.ProductInstance) (model.ProductInstance, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	INSERT INTO PRODUCT_INSTANCE(id, product_id, market_id, price,  created_at) 
		values (default, ?, ?, ?, default)
	RETURNING *
	`, product.Product.Id, product.Market.Id, product.Price)

	_, err := statement.LoadContext(ctx, &product)
	if err != nil {
		return model.ProductInstance{}, util.MakeErrorUnknown(err)
	}

	return product, nil
}

func (p ProductInstance) GetLastProductInstanceByProductId(ctx context.Context, id int64) (model.ProductInstance, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`SELECT * FROM PRODUCT_INSTANCE
         where product_id = ?
        ORDER BY created_at DESC
        LIMIT  1`, id)

	var product model.ProductInstance
	err := statement.LoadOne(&product)

	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return model.ProductInstance{}, util.MakeError(util.NOT_FOUND, fmt.Sprintf("Product Instance %d not found", id))
		}
		return model.ProductInstance{}, util.MakeErrorUnknown(err)
	}

	return product, nil
}

func (p ProductInstance) GetProductInstanceByPurchase(ctx context.Context, purchaseId int64) ([]model.ProductInstance, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	SELECT * FROM product_instance pi,
              product p,
              purchase_item pui
        WHERE 1=1
            AND pi.id = pui.product_instance_id
            AND p.id = pi.product_id
            AND pui.purchase_id = ?
	`, purchaseId)

	var product []model.ProductInstance
	_, err := statement.LoadContext(ctx, &product)
	if err != nil {
		return []model.ProductInstance{}, util.MakeErrorUnknown(err)
	}

	return product, nil
}
