package ppob

import (
	"context"

	"github.com/google/uuid"
)

// Reference gold price (e.g. Pegadaian Antam gold ~ Rp 1.450.000 / gram)
const LiveGoldRatePerGram = 1450000.0

type Service interface {
	PayBill(ctx context.Context, userID uuid.UUID, dto *PayBillDTO) (*PPOBTransaction, error)
	GetGoldAccount(ctx context.Context, userID uuid.UUID) (*GoldSavingsAccount, error)
	ConvertBalanceToGold(ctx context.Context, userID uuid.UUID, amountIDR float64) (*GoldSavingsAccount, error)
	CreateDonation(ctx context.Context, userID uuid.UUID, dto *CreateDonationDTO) (*WasteDonation, error)
	GetPPOBHistory(ctx context.Context, userID uuid.UUID, limit int) ([]PPOBTransaction, error)
	GetDonationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]WasteDonation, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) PayBill(ctx context.Context, userID uuid.UUID, dto *PayBillDTO) (*PPOBTransaction, error) {
	adminFee := 1500.0 // Standard admin fee Rp 1.500
	total := dto.Amount + adminFee

	bill := &PPOBTransaction{
		UserID:        userID,
		Type:          dto.Type,
		TargetNumber:  dto.TargetNumber,
		Amount:        dto.Amount,
		AdminFee:      adminFee,
		TotalDeducted: total,
	}

	if err := s.repo.PayBillWithWallet(ctx, userID, bill); err != nil {
		return nil, err
	}

	return bill, nil
}

func (s *service) GetGoldAccount(ctx context.Context, userID uuid.UUID) (*GoldSavingsAccount, error) {
	return s.repo.GetOrCreateGoldAccount(ctx, userID)
}

func (s *service) ConvertBalanceToGold(ctx context.Context, userID uuid.UUID, amountIDR float64) (*GoldSavingsAccount, error) {
	return s.repo.ConvertBalanceToGold(ctx, userID, amountIDR, LiveGoldRatePerGram)
}

func (s *service) CreateDonation(ctx context.Context, userID uuid.UUID, dto *CreateDonationDTO) (*WasteDonation, error) {
	donation := &WasteDonation{
		UserID:                 userID,
		CharityInstitutionName: dto.CharityInstitutionName,
		Amount:                 dto.Amount,
		Notes:                  dto.Notes,
		IsAnonymous:            dto.IsAnonymous,
	}

	if err := s.repo.DonateFromWallet(ctx, donation); err != nil {
		return nil, err
	}

	return donation, nil
}

func (s *service) GetPPOBHistory(ctx context.Context, userID uuid.UUID, limit int) ([]PPOBTransaction, error) {
	return s.repo.GetPPOBHistory(ctx, userID, limit)
}

func (s *service) GetDonationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]WasteDonation, error) {
	return s.repo.GetDonationHistory(ctx, userID, limit)
}
