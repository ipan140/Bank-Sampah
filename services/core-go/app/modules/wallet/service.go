package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	GetWallet(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	CreditBalance(ctx context.Context, userID uuid.UUID, amount float64, desc string, refID uuid.UUID) (*LedgerTransaction, error)
	RequestRedeem(ctx context.Context, userID uuid.UUID, dto *CreateRedeemDTO) (*RedeemRequest, error)
	ProcessRedeem(ctx context.Context, redeemID uuid.UUID, dto *ProcessRedeemDTO) (*RedeemRequest, error)
	GetHistory(ctx context.Context, userID uuid.UUID, page, limit int) ([]LedgerTransaction, int64, error)
	ListRedeems(ctx context.Context, page, limit int, status string) ([]RedeemRequest, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetWallet(ctx context.Context, userID uuid.UUID) (*Wallet, error) {
	return s.repo.GetWalletByUserID(ctx, userID)
}

func (s *service) CreditBalance(ctx context.Context, userID uuid.UUID, amount float64, desc string, refID uuid.UUID) (*LedgerTransaction, error) {
	var refPtr *uuid.UUID
	if refID != uuid.Nil {
		refPtr = &refID
	}
	return s.repo.CreditBalance(ctx, userID, amount, desc, refPtr)
}

func (s *service) RequestRedeem(ctx context.Context, userID uuid.UUID, dto *CreateRedeemDTO) (*RedeemRequest, error) {
	if dto.Amount < 10000 {
		return nil, errors.New("minimum redeem amount is Rp 10.000")
	}

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet.Balance < dto.Amount {
		return nil, errors.New("insufficient balance for redeem")
	}

	redeem := &RedeemRequest{
		UserID:        userID,
		Amount:        dto.Amount,
		PaymentMethod: dto.PaymentMethod,
		AccountName:   dto.AccountName,
		AccountNumber: dto.AccountNumber,
		Status:        RedeemPending,
	}

	if err := s.repo.CreateRedeem(ctx, redeem); err != nil {
		return nil, err
	}

	// Hold/Debit the balance immediately
	desc := fmt.Sprintf("Penarikan Saldo (%s %s)", dto.PaymentMethod, dto.AccountNumber)
	refID := redeem.ID
	if _, err := s.repo.DebitBalance(ctx, userID, dto.Amount, desc, &refID); err != nil {
		return nil, err
	}

	return redeem, nil
}

func (s *service) ProcessRedeem(ctx context.Context, redeemID uuid.UUID, dto *ProcessRedeemDTO) (*RedeemRequest, error) {
	redeem, err := s.repo.GetRedeemByID(ctx, redeemID)
	if err != nil {
		return nil, errors.New("redeem request not found")
	}

	if redeem.Status != RedeemPending {
		return nil, errors.New("redeem request is already processed")
	}

	redeem.Status = dto.Status
	redeem.Notes = dto.Notes

	// If rejected, refund the money back to the user
	if dto.Status == RedeemRejected {
		desc := fmt.Sprintf("Pengembalian dana penarikan ditolak: %s", dto.Notes)
		refID := redeem.ID
		_, _ = s.repo.CreditBalance(ctx, redeem.UserID, redeem.Amount, desc, &refID)
	}

	if err := s.repo.UpdateRedeem(ctx, redeem); err != nil {
		return nil, err
	}

	return redeem, nil
}

func (s *service) GetHistory(ctx context.Context, userID uuid.UUID, page, limit int) ([]LedgerTransaction, int64, error) {
	return s.repo.GetTransactions(ctx, userID, page, limit)
}

func (s *service) ListRedeems(ctx context.Context, page, limit int, status string) ([]RedeemRequest, int64, error) {
	return s.repo.ListRedeems(ctx, page, limit, status)
}
