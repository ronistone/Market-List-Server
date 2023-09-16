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
	purchase, err := p.PurchaseRepository.GetPurchaseById(ctx, purchaseId)
	if err != nil {
		return model.Purchase{}, err
	}

	if purchaseItem.Quantity == 0 {
		purchaseItem.Quantity = 1
	}

	purchaseItem.ProductInstance.Market = &purchase.Market

	product, err := p.processProduct(ctx, purchaseItem)
	if err != nil {
		return model.Purchase{}, err
	}
	purchaseItem.ProductInstance.Product = product

	productInstance, err := p.processProductInstance(ctx, purchaseItem)
	if err != nil {
		return model.Purchase{}, err
	}
	purchaseItem.ProductInstance.Id = productInstance.Id

	_, err = p.PurchaseRepository.AddPurchaseItem(ctx, purchaseId, purchaseItem)
	if err != nil {
		return model.Purchase{}, err
	}

	return p.GetPurchase(ctx, purchaseId)
}

func (p Purchase) processProduct(ctx context.Context, purchaseItem model.PurchaseItem) (model.Product, error) {
	var productFound *model.Product
	var product = purchaseItem.ProductInstance.Product

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
	}
	return *productFound, nil
}

func (p Purchase) processProductInstance(ctx context.Context, purchaseItem model.PurchaseItem) (model.ProductInstance, error) {
	var (
		productInstanceFound *model.ProductInstance
		productInstance      = purchaseItem.ProductInstance
	)
	if productInstance.Product.Id != nil {
		productInstanceQuery, err := p.ProductService.GetLastProductInstanceByProductId(ctx, *productInstance.Product.Id)
		if err == nil {
			productInstanceFound = &productInstanceQuery
		}
	} else {
		return model.ProductInstance{}, util.MakeError(util.UNKNOWN_ERROR, "Without Product Id cant create product instance")
	}
	if productInstanceFound == nil || p.calculateRealPrice(*productInstanceFound) != p.calculateRealPrice(productInstance) {
		createdInstance, err := p.ProductService.CreateInstance(ctx, productInstance)
		if err != nil {
			return model.ProductInstance{}, err
		}
		productInstanceFound = &createdInstance
	}

	return *productInstanceFound, nil
}

func (p Purchase) calculateRealPrice(productInstance model.ProductInstance) float64 {
	return float64(productInstance.Price) / math.Pow10(2)
}

func (p Purchase) RemoveItem(ctx context.Context, purchaseId int64, purchaseItemId int64) (model.Purchase, error) {
	_, err := p.PurchaseRepository.RemovePurchaseItem(ctx, purchaseId, purchaseItemId)
	if err != nil {
		return model.Purchase{}, err
	}
	return p.GetPurchase(ctx, purchaseId)
}

func (p Purchase) UpdateItem(ctx context.Context, purchaseId int64, purchaseItemId int64, item model.PurchaseItem) (model.Purchase, error) {
	purchaseItem, err := p.PurchaseRepository.GetPurchaseItemById(ctx, purchaseItemId)
	if err != nil {
		return model.Purchase{}, util.MakeError(util.NOT_FOUND, "Failed to get purchase Item")
	}

	purchaseItem.Purchased = item.Purchased
	purchaseItem.Quantity = item.Quantity

	err = p.PurchaseRepository.UpdatePurchaseItem(ctx, purchaseId, purchaseItemId, purchaseItem)
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

	var instance model.ProductInstance
	for _, item := range purchase.Items {
		instance = item.ProductInstance

		if item.Purchased {
			purchase.TotalSpent += instance.Price * int64(item.Quantity)
		}

		purchase.TotalExpected += instance.Price * int64(item.Quantity)
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
