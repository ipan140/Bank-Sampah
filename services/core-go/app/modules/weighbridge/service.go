package weighbridge

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type SettlementCallback func(ctx context.Context, nasabahID uuid.UUID, amount float64, description string, referenceID uuid.UUID) error

type Service interface {
	CreateTicket(ctx context.Context, operatorID uuid.UUID, dto *CreateTicketDTO) (*ScaleTicket, error)
	GetTicketByID(ctx context.Context, id uuid.UUID) (*ScaleTicket, error)
	ListTickets(ctx context.Context, filter ScaleTicketFilter) ([]ScaleTicket, int64, error)
	SetSettlementCallback(cb SettlementCallback)
}

type service struct {
	repo               Repository
	settlementCallback SettlementCallback
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SetSettlementCallback(cb SettlementCallback) {
	s.settlementCallback = cb
}

func (s *service) CreateTicket(ctx context.Context, operatorID uuid.UUID, dto *CreateTicketDTO) (*ScaleTicket, error) {
	if dto.GrossWeight <= dto.TareWeight {
		return nil, errors.New("gross weight must be greater than tare weight")
	}

	netWeight := dto.GrossWeight - dto.TareWeight
	totalPrice := netWeight * dto.PricePerKg

	ticket := &ScaleTicket{
		PickupID:    dto.PickupID,
		NasabahID:   dto.NasabahID,
		OperatorID:  operatorID,
		Category:    dto.Category,
		GrossWeight: dto.GrossWeight,
		TareWeight:  dto.TareWeight,
		NetWeight:   netWeight,
		PricePerKg:  dto.PricePerKg,
		TotalPrice:  totalPrice,
		DeviceID:    dto.DeviceID,
		IsSettled:   false,
	}

	if err := s.repo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	// Trigger wallet credit settlement if callback registered
	if s.settlementCallback != nil && totalPrice > 0 {
		desc := "Penimbangan Sampah Kategori " + string(ticket.Category)
		if err := s.settlementCallback(ctx, ticket.NasabahID, ticket.TotalPrice, desc, ticket.ID); err == nil {
			_ = s.repo.MarkSettled(ctx, ticket.ID)
			ticket.IsSettled = true
		}
	}

	return ticket, nil
}

func (s *service) GetTicketByID(ctx context.Context, id uuid.UUID) (*ScaleTicket, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) ListTickets(ctx context.Context, filter ScaleTicketFilter) ([]ScaleTicket, int64, error) {
	return s.repo.FindAll(ctx, filter)
}
