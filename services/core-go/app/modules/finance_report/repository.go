package finance_report

import (
	"context"

	"bank-sampah/core-go/app/modules/inventory"
	"bank-sampah/core-go/app/modules/wallet"
	"bank-sampah/core-go/app/modules/weighbridge"

	"gorm.io/gorm"
)

type Repository interface {
	CreateExpense(ctx context.Context, expense *OperationalExpense) error
	GetExpenses(ctx context.Context, limit int) ([]OperationalExpense, error)
	GetProfitAndLoss(ctx context.Context) (*ProfitAndLossStatement, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateExpense(ctx context.Context, expense *OperationalExpense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *repository) GetExpenses(ctx context.Context, limit int) ([]OperationalExpense, error) {
	var list []OperationalExpense
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	err := r.db.WithContext(ctx).Order("expense_date DESC").Limit(limit).Find(&list).Error
	return list, err
}

func (r *repository) GetProfitAndLoss(ctx context.Context) (*ProfitAndLossStatement, error) {
	pnl := &ProfitAndLossStatement{}

	// 1. Total B2B Revenue from Sales Orders
	var totalRev float64
	r.db.WithContext(ctx).Model(&inventory.B2BSaleOrder{}).Select("COALESCE(SUM(total_revenue), 0)").Scan(&totalRev)
	pnl.TotalB2BRevenue = totalRev

	// 2. Total COGS from Scale Tickets paid to Nasabah
	var totalCOGS float64
	r.db.WithContext(ctx).Model(&weighbridge.ScaleTicket{}).Select("COALESCE(SUM(total_price), 0)").Scan(&totalCOGS)
	pnl.TotalCOGS = totalCOGS

	// 3. Gross Profit
	pnl.GrossProfit = pnl.TotalB2BRevenue - pnl.TotalCOGS
	if pnl.TotalB2BRevenue > 0 {
		pnl.GrossProfitMarginPct = (pnl.GrossProfit / pnl.TotalB2BRevenue) * 100
	}

	// 4. Operational Expenses
	var totalExp float64
	r.db.WithContext(ctx).Model(&OperationalExpense{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalExp)
	pnl.TotalExpenses = totalExp

	// 5. Net Profit
	pnl.NetProfit = pnl.GrossProfit - pnl.TotalExpenses

	// 6. Total Nasabah Payouts (completed redeems)
	var totalPayouts float64
	r.db.WithContext(ctx).Model(&wallet.RedeemRequest{}).
		Where("status = ?", wallet.RedeemCompleted).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPayouts)
	pnl.TotalNasabahPayouts = totalPayouts

	return pnl, nil
}
