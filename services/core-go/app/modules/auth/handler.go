package auth

import (
	"net/http"
	"strings"

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

// Register godoc
// @Summary Register a new User
// @Description Register a new User (Nasabah, Kurir, Mitra) in the system
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} AuthResponse
// @Param request body RegisterRequest true "Payload"
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return response.BadRequest(c, "Name, email, and password are required", nil)
	}

	res, err := h.service.Register(c.Request().Context(), &req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "User registered successfully", res)
}

// Login godoc
// @Summary User Login
// @Description Authenticate user credentials and return JWT bearer token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} AuthResponse
// @Param request body LoginRequest true "Credentials"
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	res, err := h.service.Login(c.Request().Context(), &req)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Login successful", res)
}

func (h *Handler) Logout(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 {
		_ = h.service.Logout(c.Request().Context(), parts[1])
	}

	return response.Success(c, http.StatusOK, "Logout successful", nil)
}

func (h *Handler) Profile(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	user, err := h.service.GetProfile(c.Request().Context(), uid)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}
