package service

import (
	"context"
	"github.com/ronistone/market-list/src/model"
	"github.com/ronistone/market-list/src/repository"
	"github.com/ronistone/market-list/src/util"
	"math"
)

type PurchaseService interface {
	CreatePurchase(ctx context.Context, purchase model.Purchase) (model.Purchase, error)
	AddItem(ctx context.Context, purchaseId int64, purchaseItem model.PurchaseItem) (model.Purchase, error)
	RemoveItem(ctx context.Context, purchaseId int64, purchaseItemId int64) (model.Purchase, error)
	UpdateItem(ctx context.Context, purchaseId int64, purchaseItemId int64, item model.PurchaseItem) (model.Purchase, error)
	GetPurchase(ctx context.Context, id int64) (model.Purchase, error)
	GetAllPurchase(ctx context.Context) ([]model.Purchase, error)
	DeletePurchase(ctx context.Context, id int64) error
	GetItem(ctx context.Context, purchaseId int64, purchaseItemId int64) (model.PurchaseItem, error)
}
type Purchase struct {
	PurchaseRepository repository.PurchaseRepository
	ProductService     ProductService
}

func CreatePurchaseService(
	purchaseRepository repository.PurchaseRepository,
	productService ProductService,
) PurchaseService {
	return &Purchase{
		PurchaseRepository: purchaseRepository,
		ProductService:     productService,
	}
}

func (p Purchase) CreatePurchase(ctx context.Context, purchase model.Purchase) (model.Purchase, error) {
	if purchase.MarketId == nil || purchase.UserId == nil {
		return model.Purchase{}, util.MakeError(util.INVALID_INPUT, "Invalid marketId or userId")
	}

	purchase.Market.Id = purchase.MarketId
	purchase.User.Id = purchase.UserId

	created, err := p.PurchaseRepository.CreatePurchase(ctx, purchase)
	if err != nil {
		return model.Purchase{}, err
	}
	return p.GetPurchase(ctx, *created.Id)
}

func (p Purchase) AddItem(ctx context.Context, purchaseId int64, purchaseItem model.PurchaseItem) (model.Purchase, error) {
	_, err := p.PurchaseRepository.GetPurchaseById(ctx, purchaseId)
	if err != nil {
		return model.Purchase{}, err
	}

	if purchaseItem.Quantity == 0 {
		purchaseItem.Quantity = 1
	}

	product, err := p.processProduct(ctx, purchaseItem)
	if err != nil {
		return model.Purchase{}, err
	}
	purchaseItem.Product = product

	_, err = p.PurchaseRepository.AddPurchaseItem(ctx, purchaseId, purchaseItem)
	if err != nil {
		return model.Purchase{}, err
	}

	return p.GetPurchase(ctx, purchaseId)
}

func (p Purchase) processProduct(ctx context.Context, purchaseItem model.PurchaseItem) (model.Product, error) {
	var productFound *model.Product
	var product = purchaseItem.Product

	if product.Id != nil {
		productQuery, err := p.ProductService.GetById(ctx, *product.Id)
		if err == nil {
			productFound = &productQuery
		}
	}

	if product.Ean != nil && len(*product.Ean) > 0 {
		productQuery, err := p.ProductService.GetByEan(ctx, *product.Ean)
		if err == nil {
			productFound = &productQuery
		}

	}

	if productFound == nil {
		createdProduct, err := p.ProductService.Create(ctx, product)
		if err != nil {
			return model.Product{}, err
		}
		productFound = &createdProduct
	} else {
		productFound.Name = product.Name
		productFound.Ean = product.Ean
		productFound.Size = product.Size
		productFound.Unit = product.Unit
		updatedProduct, err := p.ProductService.Update(ctx, *productFound)
		if err != nil {
			return model.Product{}, err
		}
		productFound = &updatedProduct
	}
	return *productFound, nil
}

func (p Purchase) calculateRealPrice(productInstance model.PurchaseItem) *float64 {
	if productInstance.Price != nil {
		price := float64(*productInstance.Price) / math.Pow10(2)
		return &price
	}
	return nil
}

func (p Purchase) RemoveItem(ctx context.Context, purchaseId int64, purchaseItemId int64) (model.Purchase, error) {
	_, err := p.PurchaseRepository.RemovePurchaseItem(ctx, purchaseId, purchaseItemId)
	if err != nil {
		return model.Purchase{}, err
	}
	return p.GetPurchase(ctx, purchaseId)
}

func (p Purchase) UpdateItem(ctx context.Context, purchaseId int64, purchaseItemId int64, item model.PurchaseItem) (model.Purchase, error) {
	_, err := p.PurchaseRepository.GetPurchaseItemById(ctx, purchaseId, purchaseItemId)
	if err != nil {
		return model.Purchase{}, util.MakeError(util.NOT_FOUND, "Failed to get purchase Item")
	}

	product, err := p.processProduct(ctx, item)
	if err != nil {
		return model.Purchase{}, err
	}
	item.Product = product

	err = p.PurchaseRepository.UpdatePurchaseItem(ctx, purchaseId, purchaseItemId, item)
	if err != nil {
		return model.Purchase{}, err
	}
	return p.GetPurchase(ctx, purchaseId)
}

func (p Purchase) GetPurchase(ctx context.Context, id int64) (model.Purchase, error) {
	purchase, err := p.PurchaseRepository.GetPurchaseByIdFetchItems(ctx, id)
	if err != nil {
		return model.Purchase{}, err
	}

	purchase.TotalSpent = 0
	purchase.TotalExpected = 0

	for _, item := range purchase.Items {

		if item.Price == nil {
			continue
		}

		if item.Purchased {
			purchase.TotalSpent += *item.Price * int64(item.Quantity)
		}

		purchase.TotalExpected += *item.Price * int64(item.Quantity)
	}

	return purchase, nil
}

func (p Purchase) GetAllPurchase(ctx context.Context) ([]model.Purchase, error) {
	var userId int64 = 1 // TODO user hardcoded
	return p.PurchaseRepository.ListPurchase(ctx, userId)
}

func (p Purchase) DeletePurchase(ctx context.Context, id int64) error {
	purchase, err := p.GetPurchase(ctx, id)
	if err != nil {
		return err
	}

	for _, item := range purchase.Items {
		_, err := p.RemoveItem(ctx, *purchase.Id, *item.Id)
		if err != nil {
			return err
		}
	}

	return p.PurchaseRepository.DeletePurchase(ctx, id)
}

func (p Purchase) GetItem(ctx context.Context, purchaseId int64, purchaseItemId int64) (model.PurchaseItem, error) {
	return p.PurchaseRepository.GetPurchaseItemById(ctx, purchaseId, purchaseItemId)
}
