package auth

import (
	"context"
	"errors"
	"time"

	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	Logout(ctx context.Context, token string) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
}

type service struct {
	repo      Repository
	cacheRepo CacheRepository
	cfg       *config.Config
}

func NewService(repo Repository, cacheRepo CacheRepository, cfg *config.Config) Service {
	return &service{
		repo:      repo,
		cacheRepo: cacheRepo,
		cfg:       cfg,
	}
}

func (s *service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	existing, _ := s.repo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = RoleNasabah
	}

	user := &User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     role,
		Address:  req.Address,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := middleware.GenerateToken(user.ID.String(), user.Email, string(user.Role), s.cfg)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := middleware.GenerateToken(user.ID.String(), user.Email, string(user.Role), s.cfg)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	return s.cacheRepo.SetTokenBlacklist(ctx, token, 24*time.Hour)
}

func (s *service) GetProfile(ctx context.Context, userID uuid.UUID) (*User, error) {
	return s.repo.FindByID(ctx, userID)
}
