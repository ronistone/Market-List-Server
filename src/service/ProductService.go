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
	GetProductsByPurchase(ctx context.Context, purchaseId int64) ([]model.ProductInstance, error)
	GetLastProductInstanceByProductId(ctx context.Context, id int64) (model.ProductInstance, error)
	CreateInstance(ctx context.Context, instance model.ProductInstance) (model.ProductInstance, error)
}

type Product struct {
	ProductRepository         repository.ProductRepository
	ProductInstanceRepository repository.ProductInstanceRepository
}

func CreateProductService(productRepository repository.ProductRepository, productInstanceRepository repository.ProductInstanceRepository) ProductService {
	return &Product{
		ProductRepository:         productRepository,
		ProductInstanceRepository: productInstanceRepository,
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

func (p Product) GetProductsByPurchase(ctx context.Context, purchaseId int64) ([]model.ProductInstance, error) {
	return p.ProductInstanceRepository.GetProductInstanceByPurchase(ctx, purchaseId)
}

func (p Product) GetLastProductInstanceByProductId(ctx context.Context, productId int64) (model.ProductInstance, error) {
	return p.ProductInstanceRepository.GetLastProductInstanceByProductId(ctx, productId)
}

func (p Product) CreateInstance(ctx context.Context, instance model.ProductInstance) (model.ProductInstance, error) {
	return p.ProductInstanceRepository.CreateProduct(ctx, instance)
}
