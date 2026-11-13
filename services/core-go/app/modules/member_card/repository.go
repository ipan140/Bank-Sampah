package member_card

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/app/modules/wallet"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetOrCreateCard(ctx context.Context, userID uuid.UUID) (*MemberCard, error)
	FindByQRCode(ctx context.Context, qrToken string) (*MemberCard, error)
	GetMemberInfo(ctx context.Context, qrToken string) (*MemberInfoResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetOrCreateCard(ctx context.Context, userID uuid.UUID) (*MemberCard, error) {
	var card MemberCard
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&card).Error
	if err == nil {
		return &card, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		cardNum := fmt.Sprintf("BS-%d-%04d", time.Now().Format("060102"), rand.Intn(9999))
		qrToken := fmt.Sprintf("BS-QR-%s-%s", userID.String()[:8], uuid.New().String()[:8])
		card = MemberCard{
			ID:          uuid.New(),
			UserID:      userID,
			CardNumber:  cardNum,
			QRCodeToken: qrToken,
			IsActive:    true,
			IssuedAt:    time.Now(),
		}
		if err := r.db.WithContext(ctx).Create(&card).Error; err != nil {
			return nil, err
		}
		return &card, nil
	}
	return nil, err
}

func (r *repository) FindByQRCode(ctx context.Context, qrToken string) (*MemberCard, error) {
	var card MemberCard
	err := r.db.WithContext(ctx).Where("qr_code_token = ? AND is_active = ?", qrToken, true).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *repository) GetMemberInfo(ctx context.Context, qrToken string) (*MemberInfoResponse, error) {
	card, err := r.FindByQRCode(ctx, qrToken)
	if err != nil {
		return nil, errors.New("invalid or inactive member QR code")
	}

	var user auth.User
	if err := r.db.WithContext(ctx).Where("id = ?", card.UserID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var userWallet wallet.Wallet
	_ = r.db.WithContext(ctx).Where("user_id = ?", card.UserID).First(&userWallet)

	return &MemberInfoResponse{
		UserID:     user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		CardNumber: card.CardNumber,
		Balance:    userWallet.Balance,
		Points:     userWallet.Points,
	}, nil
}
