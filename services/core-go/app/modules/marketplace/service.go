package marketplace

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateProduct(ctx context.Context, dto *CreateProductDTO) (*EcoProduct, error)
	ListProducts(ctx context.Context, category string) ([]EcoProduct, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*EcoProduct, error)
	Checkout(ctx context.Context, userID uuid.UUID, dto *CheckoutOrderDTO) (*EcoOrder, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]EcoOrder, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateProduct(ctx context.Context, dto *CreateProductDTO) (*EcoProduct, error) {
	unit := dto.Unit
	if unit == "" {
		unit = "pcs"
	}

	product := &EcoProduct{
		Name:        dto.Name,
		Description: dto.Description,
		Category:    dto.Category,
		PriceIDR:    dto.PriceIDR,
		Stock:       dto.Stock,
		Unit:        unit,
		ImageURL:    dto.ImageURL,
		IsActive:    true,
	}

	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *service) ListProducts(ctx context.Context, category string) ([]EcoProduct, error) {
	return s.repo.GetProducts(ctx, category)
}

func (s *service) GetProductByID(ctx context.Context, id uuid.UUID) (*EcoProduct, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *service) Checkout(ctx context.Context, userID uuid.UUID, dto *CheckoutOrderDTO) (*EcoOrder, error) {
	return s.repo.CheckoutOrder(ctx, userID, dto.ShippingAddress, dto.Items)
}

func (s *service) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]EcoOrder, error) {
	return s.repo.GetUserOrders(ctx, userID)
}
