package inventory

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"bank-sampah/core-go/app/modules/weighbridge"

	"github.com/google/uuid"
)

type Service interface {
	GetStocks(ctx context.Context) ([]WasteStock, error)
	RecordStockIn(ctx context.Context, category weighbridge.WasteCategory, quantityKg float64, pricePerKg float64, refID uuid.UUID) error
	CreateB2BSale(ctx context.Context, dto *CreateB2BSaleDTO) (*B2BSaleOrder, error)
	ListB2BSales(ctx context.Context, page, limit int) ([]B2BSaleOrder, int64, error)
	GetMutations(ctx context.Context, category string, limit int) ([]StockMutation, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetStocks(ctx context.Context) ([]WasteStock, error) {
	return s.repo.GetStocks(ctx)
}

func (s *service) RecordStockIn(ctx context.Context, category weighbridge.WasteCategory, quantityKg float64, pricePerKg float64, refID uuid.UUID) error {
	var refPtr *uuid.UUID
	if refID != uuid.Nil {
		refPtr = &refID
	}
	notes := fmt.Sprintf("Setoran timbangan %.2f kg @ Rp %.0f", quantityKg, pricePerKg)
	return s.repo.AddStock(ctx, category, quantityKg, pricePerKg, refPtr, notes)
}

func (s *service) CreateB2BSale(ctx context.Context, dto *CreateB2BSaleDTO) (*B2BSaleOrder, error) {
	invNum := fmt.Sprintf("INV-B2B-%d-%04d", time.Now().Unix(), rand.Intn(9999))
	sale := &B2BSaleOrder{
		BuyerFactory:   dto.BuyerFactory,
		InvoiceNumber:  invNum,
		Category:       dto.Category,
		WeightSoldKg:   dto.WeightSoldKg,
		SalePricePerKg: dto.SalePricePerKg,
		SaleDate:       time.Now(),
	}

	if err := s.repo.CreateB2BSale(ctx, sale); err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *service) ListB2BSales(ctx context.Context, page, limit int) ([]B2BSaleOrder, int64, error) {
	return s.repo.GetB2BSales(ctx, page, limit)
}

func (s *service) GetMutations(ctx context.Context, category string, limit int) ([]StockMutation, error) {
	return s.repo.GetMutations(ctx, category, limit)
}
