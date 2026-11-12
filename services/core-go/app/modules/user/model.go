package user

import (
	"time"

	"bank-sampah/core-go/app/modules/auth"
	"github.com/google/uuid"
)

type UserDetail struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      auth.Role `json:"role"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type AssignRoleRequest struct {
	Role auth.Role `json:"role" validate:"required"`
}

type UserFilter struct {
	Role   string `query:"role"`
	Search string `query:"search"`
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
}
