package compliance

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateCertificate(ctx context.Context, cert *EPRCertificate) error
	GetCertificateByHash(ctx context.Context, hash string) (*EPRCertificate, error)
	ListCertificates(ctx context.Context, limit int) ([]EPRCertificate, error)
	CreateSubscription(ctx context.Context, sub *B2BSubscriptionContract) error
	ListSubscriptions(ctx context.Context) ([]B2BSubscriptionContract, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCertificate(ctx context.Context, cert *EPRCertificate) error {
	return r.db.WithContext(ctx).Create(cert).Error
}

func (r *repository) GetCertificateByHash(ctx context.Context, hash string) (*EPRCertificate, error) {
	var cert EPRCertificate
	err := r.db.WithContext(ctx).Where("verification_hash = ?", hash).First(&cert).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *repository) ListCertificates(ctx context.Context, limit int) ([]EPRCertificate, error) {
	var list []EPRCertificate
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	err := r.db.WithContext(ctx).Order("issued_date DESC").Limit(limit).Find(&list).Error
	return list, err
}

func (r *repository) CreateSubscription(ctx context.Context, sub *B2BSubscriptionContract) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *repository) ListSubscriptions(ctx context.Context) ([]B2BSubscriptionContract, error) {
	var list []B2BSubscriptionContract
	err := r.db.WithContext(ctx).Order("start_date DESC").Find(&list).Error
	return list, err
}
