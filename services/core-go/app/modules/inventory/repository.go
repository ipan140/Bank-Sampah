package inventory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bank-sampah/core-go/app/modules/weighbridge"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetStocks(ctx context.Context) ([]WasteStock, error)
	GetStockByCategory(ctx context.Context, category weighbridge.WasteCategory) (*WasteStock, error)
	AddStock(ctx context.Context, category weighbridge.WasteCategory, quantityKg float64, costPerKg float64, refID *uuid.UUID, notes string) error
	CreateB2BSale(ctx context.Context, sale *B2BSaleOrder) error
	GetB2BSales(ctx context.Context, page, limit int) ([]B2BSaleOrder, int64, error)
	GetMutations(ctx context.Context, category string, limit int) ([]StockMutation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetStocks(ctx context.Context) ([]WasteStock, error) {
	var stocks []WasteStock
	err := r.db.WithContext(ctx).Order("category ASC").Find(&stocks).Error
	return stocks, err
}

func (r *repository) GetStockByCategory(ctx context.Context, category weighbridge.WasteCategory) (*WasteStock, error) {
	var stock WasteStock
	err := r.db.WithContext(ctx).Where("category = ?", category).First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			stock = WasteStock{
				ID:           uuid.New(),
				Category:     category,
				CurrentKg:    0,
				AvgCostPerKg: 0,
				LastUpdated:  time.Now(),
			}
			_ = r.db.WithContext(ctx).Create(&stock)
			return &stock, nil
		}
		return nil, err
	}
	return &stock, nil
}

func (r *repository) AddStock(ctx context.Context, category weighbridge.WasteCategory, quantityKg float64, costPerKg float64, refID *uuid.UUID, notes string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stock WasteStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("category = ?", category).First(&stock).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				stock = WasteStock{
					ID:           uuid.New(),
					Category:     category,
					CurrentKg:    0,
					AvgCostPerKg: costPerKg,
					LastUpdated:  time.Now(),
				}
				if err := tx.Create(&stock).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		stockBefore := stock.CurrentKg
		stockAfter := stockBefore + quantityKg

		// Weighted average cost per kg
		if stockAfter > 0 {
			stock.AvgCostPerKg = ((stockBefore * stock.AvgCostPerKg) + (quantityKg * costPerKg)) / stockAfter
		}
		stock.CurrentKg = stockAfter
		stock.LastUpdated = time.Now()

		if err := tx.Save(&stock).Error; err != nil {
			return err
		}

		mutation := StockMutation{
			Category:    category,
			Type:        MovementInWeighbridge,
			QuantityKg:  quantityKg,
			StockBefore: stockBefore,
			StockAfter:  stockAfter,
			ReferenceID: refID,
			Notes:       notes,
			CreatedAt:   time.Now(),
		}
		return tx.Create(&mutation).Error
	})
}

func (r *repository) CreateB2BSale(ctx context.Context, sale *B2BSaleOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stock WasteStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("category = ?", sale.Category).First(&stock).Error; err != nil {
			return errors.New("waste stock category not found")
		}

		if stock.CurrentKg < sale.WeightSoldKg {
			return fmt.Errorf("insufficient stock: available %.2f kg, requested %.2f kg", stock.CurrentKg, sale.WeightSoldKg)
		}

		stockBefore := stock.CurrentKg
		stockAfter := stockBefore - sale.WeightSoldKg
		stock.CurrentKg = stockAfter
		stock.LastUpdated = time.Now()

		// Set cost price from inventory average cost
		sale.CostPricePerKg = stock.AvgCostPerKg
		sale.TotalCost = sale.WeightSoldKg * sale.CostPricePerKg
		sale.TotalRevenue = sale.WeightSoldKg * sale.SalePricePerKg
		sale.GrossProfit = sale.TotalRevenue - sale.TotalCost

		if err := tx.Save(&stock).Error; err != nil {
			return err
		}

		if err := tx.Create(sale).Error; err != nil {
			return err
		}

		mutation := StockMutation{
			Category:    sale.Category,
			Type:        MovementOutSale,
			QuantityKg:  sale.WeightSoldKg,
			StockBefore: stockBefore,
			StockAfter:  stockAfter,
			ReferenceID: &sale.ID,
			Notes:       fmt.Sprintf("Penjualan ke %s (Inv: %s)", sale.BuyerFactory, sale.InvoiceNumber),
			CreatedAt:   time.Now(),
		}
		return tx.Create(&mutation).Error
	})
}

func (r *repository) GetB2BSales(ctx context.Context, page, limit int) ([]B2BSaleOrder, int64, error) {
	var sales []B2BSaleOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&B2BSaleOrder{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("sale_date DESC").Limit(limit).Offset(offset).Find(&sales).Error
	return sales, total, err
}

func (r *repository) GetMutations(ctx context.Context, category string, limit int) ([]StockMutation, error) {
	var mutations []StockMutation
	query := r.db.WithContext(ctx)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	err := query.Order("created_at DESC").Limit(limit).Find(&mutations).Error
	return mutations, err
}
