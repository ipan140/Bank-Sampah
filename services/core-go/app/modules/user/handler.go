package user

import (
	"net/http"

	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListUsers(c echo.Context) error {
	var filter UserFilter
	if err := c.Bind(&filter); err != nil {
		return response.BadRequest(c, "Invalid filter params", err.Error())
	}

	users, total, err := h.service.ListUsers(c.Request().Context(), filter)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	meta := response.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return response.SuccessWithMeta(c, http.StatusOK, "Users retrieved successfully", users, meta)
}

func (h *Handler) GetUser(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	user, err := h.service.GetUserByID(c.Request().Context(), uid)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "User details retrieved", user)
}

func (h *Handler) UpdateProfile(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	updated, err := h.service.UpdateProfile(c.Request().Context(), uid, &req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Profile updated successfully", updated)
}

func (h *Handler) AssignRole(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var req AssignRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	if req.Role != auth.RoleAdmin && req.Role != auth.RoleNasabah && req.Role != auth.RoleKurir && req.Role != auth.RoleMitra {
		return response.BadRequest(c, "Invalid role specified", nil)
	}

	if err := h.service.AssignRole(c.Request().Context(), uid, req.Role); err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Role assigned successfully", nil)
}
