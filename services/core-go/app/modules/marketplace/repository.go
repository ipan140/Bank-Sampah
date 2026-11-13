package marketplace

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"bank-sampah/core-go/app/modules/wallet"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateProduct(ctx context.Context, p *EcoProduct) error
	GetProducts(ctx context.Context, category string) ([]EcoProduct, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*EcoProduct, error)
	CheckoutOrder(ctx context.Context, userID uuid.UUID, shippingAddress string, items []CheckoutItemDTO) (*EcoOrder, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]EcoOrder, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateProduct(ctx context.Context, p *EcoProduct) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *repository) GetProducts(ctx context.Context, category string) ([]EcoProduct, error) {
	var list []EcoProduct
	query := r.db.WithContext(ctx).Where("is_active = ? AND stock > 0", true)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *repository) GetProductByID(ctx context.Context, id uuid.UUID) (*EcoProduct, error) {
	var p EcoProduct
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) CheckoutOrder(ctx context.Context, userID uuid.UUID, shippingAddress string, items []CheckoutItemDTO) (*EcoOrder, error) {
	var order *EcoOrder

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var orderItems []EcoOrderItem

		// Check stock and calculate total
		for _, itemReq := range items {
			var product EcoProduct
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", itemReq.ProductID).First(&product).Error; err != nil {
				return fmt.Errorf("product %s not found", itemReq.ProductID)
			}

			if product.Stock < itemReq.Quantity {
				return fmt.Errorf("insufficient stock for product %s (available: %d)", product.Name, product.Stock)
			}

			subtotal := float64(itemReq.Quantity) * product.PriceIDR
			totalAmount += subtotal

			// Deduct product stock
			product.Stock -= itemReq.Quantity
			if err := tx.Save(&product).Error; err != nil {
				return err
			}

			orderItems = append(orderItems, EcoOrderItem{
				ID:          uuid.New(),
				ProductID:   product.ID,
				ProductName: product.Name,
				Quantity:    itemReq.Quantity,
				UnitPrice:   product.PriceIDR,
				Subtotal:    subtotal,
			})
		}

		// Deduct Wallet Balance
		var userWallet wallet.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&userWallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		if userWallet.Balance < totalAmount {
			return fmt.Errorf("insufficient wallet balance (balance: Rp %.0f, required: Rp %.0f)", userWallet.Balance, totalAmount)
		}

		userWallet.Balance -= totalAmount
		if err := tx.Save(&userWallet).Error; err != nil {
			return err
		}

		// Create Eco Order
		orderNumber := fmt.Sprintf("ECO-%d-%04d", time.Now().Unix(), rand.Intn(9999))
		order = &EcoOrder{
			ID:              uuid.New(),
			OrderNumber:     orderNumber,
			UserID:          userID,
			TotalAmountIDR:  totalAmount,
			ShippingAddress: shippingAddress,
			Status:          "PAID",
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
			if err := tx.Create(&orderItems[i]).Error; err != nil {
				return err
			}
		}
		order.Items = orderItems

		// Ledger deduction
		ledger := wallet.LedgerTransaction{
			WalletID:      userWallet.ID,
			UserID:        userID,
			Type:          wallet.TxDebitRedeem,
			Amount:        totalAmount,
			BalanceBefore: userWallet.Balance + totalAmount,
			BalanceAfter:  userWallet.Balance,
			Description:   fmt.Sprintf("Pembelian Produk Daur Ulang (Order %s)", orderNumber),
			ReferenceID:   &order.ID,
			CreatedAt:     time.Now(),
		}
		return tx.Create(&ledger).Error
	})

	return order, err
}

func (r *repository) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]EcoOrder, error) {
	var orders []EcoOrder
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}
