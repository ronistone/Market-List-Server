package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/lib/pq"
	"github.com/ronistone/market-list/src/model"
	"github.com/ronistone/market-list/src/util"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)
	UpdateProduct(ctx context.Context, product model.Product) (model.Product, error)
	GetProductByEan(ctx context.Context, ean string) (model.Product, error)
	GetProductById(ctx context.Context, id int64) (model.Product, error)
	GetProductByName(ctx context.Context, name string, limit int) ([]model.Product, error)
}

type Product struct {
	DbConnection *dbr.Connection
}

func CreateProductRepository(connection *dbr.Connection) ProductRepository {
	return &Product{
		DbConnection: connection,
	}
}

func (p Product) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	INSERT INTO PRODUCT(id, ean, name, unit, size, created_at, updated_at) 
		values (default, ?, ?, ?, ?, default, default)
	RETURNING *
	`, product.Ean, product.Name, product.Unit, product.Size)

	_, err := statement.LoadContext(ctx, &product)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			return model.Product{}, util.MakeError(util.ALREADY_EXISTS, pqError.Message)
		}
		return model.Product{}, util.MakeErrorUnknown(err)
	}

	return product, nil
}

func (p Product) UpdateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	if product.Id == nil {
		return model.Product{}, util.MakeError(util.INVALID_INPUT, "invalid Product Id")
	}
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	UPDATE PRODUCT SET ean = ?, name = ?, unit = ?, size = ?, updated_at = NOW() 
		WHERE id = ?
	RETURNING *
	`, product.Ean, product.Name, product.Unit, product.Size, product.Id)

	_, err := statement.LoadContext(ctx, &product)
	if err != nil {
		return model.Product{}, util.MakeErrorUnknown(err)
	}

	return product, nil
}

func (p Product) GetProductByEan(ctx context.Context, ean string) (model.Product, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	SELECT * FROM product where ean = ?
	`, ean)

	var product model.Product
	err := statement.LoadOne(&product)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return model.Product{}, util.MakeError(util.NOT_FOUND, fmt.Sprintf("Product %s not found", ean))
		}
		return model.Product{}, util.MakeErrorUnknown(err)
	}

	return product, nil
}

func (p Product) GetProductById(ctx context.Context, id int64) (model.Product, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	SELECT * FROM product where id = ?
	`, id)

	var product model.Product
	err := statement.LoadOne(&product)

	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return model.Product{}, util.MakeError(util.NOT_FOUND, fmt.Sprintf("Product %d not found", id))
		}
		return model.Product{}, util.MakeErrorUnknown(err)
	}

	return product, nil
}

func (p Product) GetProductByName(ctx context.Context, name string, limit int) ([]model.Product, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	SELECT * FROM (
          SELECT similarity(NAME, ?) simi, * FROM PRODUCT
		  ) p
	WHERE simi > 0
	ORDER BY  simi DESC
	limit ?
	`, name, limit)

	var product []model.Product
	_, err := statement.LoadContext(ctx, &product)
	if err != nil {
		return []model.Product{}, util.MakeErrorUnknown(err)
	}

	return product, nil
}
