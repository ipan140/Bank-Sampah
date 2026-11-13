package compliance

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

type Service interface {
	IssueCertificate(ctx context.Context, dto *IssueCertificateDTO) (*EPRCertificate, error)
	VerifyCertificate(ctx context.Context, hash string) (*EPRCertificate, error)
	ListCertificates(ctx context.Context, limit int) ([]EPRCertificate, error)
	CreateSubscription(ctx context.Context, dto *CreateSubscriptionDTO) (*B2BSubscriptionContract, error)
	ListSubscriptions(ctx context.Context) ([]B2BSubscriptionContract, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) IssueCertificate(ctx context.Context, dto *IssueCertificateDTO) (*EPRCertificate, error) {
	totalKg := dto.PlasticWasteKg + dto.PaperWasteKg + dto.MetalWasteKg

	// ESG Environmental Formulas:
	// 1 kg plastic recycled ~ 1.5 kg CO2 avoided
	// 1 kg metal recycled ~ 2.0 kg CO2 avoided
	// 1 kg paper recycled ~ 0.9 kg CO2 avoided
	co2Avoided := (dto.PlasticWasteKg * 1.5) + (dto.MetalWasteKg * 2.0) + (dto.PaperWasteKg * 0.9)

	now := time.Now()
	validMonths := dto.ValidMonths
	if validMonths <= 0 {
		validMonths = 12
	}
	validUntil := now.AddDate(0, validMonths, 0)

	certNumber := fmt.Sprintf("EPR-%d-%04d", now.Unix(), rand.Intn(9999))
	rawSign := fmt.Sprintf("%s|%s|%.2f|%s", certNumber, dto.RecipientCompanyName, totalKg, now.Format(time.RFC3339))
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(rawSign)))[:24]

	cert := &EPRCertificate{
		CertificateNumber:    certNumber,
		RecipientCompanyName: dto.RecipientCompanyName,
		TaxNumberNPWP:        dto.TaxNumberNPWP,
		TotalWasteKg:         totalKg,
		PlasticWasteKg:       dto.PlasticWasteKg,
		PaperWasteKg:         dto.PaperWasteKg,
		MetalWasteKg:         dto.MetalWasteKg,
		Co2AvoidedKg:         co2Avoided,
		IssuedDate:           now,
		ValidUntil:           validUntil,
		VerificationHash:     hash,
		Status:               "VALID",
	}

	if err := s.repo.CreateCertificate(ctx, cert); err != nil {
		return nil, err
	}

	return cert, nil
}

func (s *service) VerifyCertificate(ctx context.Context, hash string) (*EPRCertificate, error) {
	return s.repo.GetCertificateByHash(ctx, hash)
}

func (s *service) ListCertificates(ctx context.Context, limit int) ([]EPRCertificate, error) {
	return s.repo.ListCertificates(ctx, limit)
}

func (s *service) CreateSubscription(ctx context.Context, dto *CreateSubscriptionDTO) (*B2BSubscriptionContract, error) {
	now := time.Now()
	duration := dto.DurationMonths
	if duration <= 0 {
		duration = 12
	}
	contractNumber := fmt.Sprintf("B2B-SUB-%d-%04d", now.Unix(), rand.Intn(9999))

	sub := &B2BSubscriptionContract{
		ContractNumber:  contractNumber,
		CompanyName:     dto.CompanyName,
		Address:         dto.Address,
		ContactPerson:   dto.ContactPerson,
		ContactPhone:    dto.ContactPhone,
		ContactEmail:    dto.ContactEmail,
		PickupFrequency: dto.PickupFrequency,
		MonthlyFeeIDR:   dto.MonthlyFeeIDR,
		StartDate:       now,
		EndDate:         now.AddDate(0, duration, 0),
		IsActive:        true,
	}

	if err := s.repo.CreateSubscription(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *service) ListSubscriptions(ctx context.Context) ([]B2BSubscriptionContract, error) {
	return s.repo.ListSubscriptions(ctx)
}
