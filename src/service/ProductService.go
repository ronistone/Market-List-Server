package service

import (
	"context"
	"github.com/ronistone/market-list/src/model"
	"github.com/ronistone/market-list/src/repository"
)

type ProductService interface {
	Create(ctx context.Context, product model.Product) (model.Product, error)
	Update(ctx context.Context, product model.Product) (model.Product, error)
	GetByName(ctx context.Context, name string) ([]model.Product, error)
	GetByEan(ctx context.Context, ean string) (model.Product, error)
	GetById(ctx context.Context, id int64) (model.Product, error)
}

type Product struct {
	ProductRepository repository.ProductRepository
}

func CreateProductService(productRepository repository.ProductRepository) ProductService {
	return &Product{
		ProductRepository: productRepository,
	}
}

func (p Product) Create(ctx context.Context, product model.Product) (model.Product, error) {
	return p.ProductRepository.CreateProduct(ctx, product)
}

func (p Product) Update(ctx context.Context, product model.Product) (model.Product, error) {
	return p.ProductRepository.UpdateProduct(ctx, product)
}

func (p Product) GetByName(ctx context.Context, name string) ([]model.Product, error) {
	return p.ProductRepository.GetProductByName(ctx, name, 5)
}

func (p Product) GetByEan(ctx context.Context, ean string) (model.Product, error) {
	return p.ProductRepository.GetProductByEan(ctx, ean)
}

func (p Product) GetById(ctx context.Context, id int64) (model.Product, error) {
	return p.ProductRepository.GetProductById(ctx, id)
}
