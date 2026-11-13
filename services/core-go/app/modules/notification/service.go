package notification

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type Service interface {
	SendNotification(ctx context.Context, dto *SendNotificationDTO) (*Notification, error)
	GetUserNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SendNotification(ctx context.Context, dto *SendNotificationDTO) (*Notification, error) {
	channel := dto.Channel
	if channel == "" {
		channel = ChannelInApp
	}

	notif := &Notification{
		UserID:        dto.UserID,
		Title:         dto.Title,
		Message:       dto.Message,
		Channel:       channel,
		Status:        StatusSent,
		ReferenceType: dto.ReferenceType,
	}

	if err := s.repo.Create(ctx, notif); err != nil {
		return nil, err
	}

	// Dispatcher Hook: WhatsApp / FCM gateway integration
	if channel == ChannelWhatsApp {
		log.Printf("[WHATSAPP GATEWAY DISPATCH] To User %s: %s - %s", dto.UserID, dto.Title, dto.Message)
	}

	return notif, nil
}

func (s *service) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]Notification, error) {
	return s.repo.GetByUserID(ctx, userID, limit)
}

func (s *service) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}
