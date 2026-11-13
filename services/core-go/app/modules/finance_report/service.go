package finance_report

import (
	"context"
	"time"
)

type Service interface {
	RecordExpense(ctx context.Context, dto *CreateExpenseDTO) (*OperationalExpense, error)
	GetExpenses(ctx context.Context, limit int) ([]OperationalExpense, error)
	GetProfitAndLoss(ctx context.Context) (*ProfitAndLossStatement, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RecordExpense(ctx context.Context, dto *CreateExpenseDTO) (*OperationalExpense, error) {
	expDate := dto.ExpenseDate
	if expDate.IsZero() {
		expDate = time.Now()
	}

	exp := &OperationalExpense{
		Category:    dto.Category,
		Amount:      dto.Amount,
		Description: dto.Description,
		ReceiptURL:  dto.ReceiptURL,
		ExpenseDate: expDate,
	}

	if err := s.repo.CreateExpense(ctx, exp); err != nil {
		return nil, err
	}

	return exp, nil
}

func (s *service) GetExpenses(ctx context.Context, limit int) ([]OperationalExpense, error) {
	return s.repo.GetExpenses(ctx, limit)
}

func (s *service) GetProfitAndLoss(ctx context.Context) (*ProfitAndLossStatement, error) {
	return s.repo.GetProfitAndLoss(ctx)
}
