package wallet

import (
	"net/http"
	"strconv"

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

// GetWallet godoc
// @Summary Get Wallet Balance
// @Description Get current logged-in user wallet balance and reward points
// @Tags wallet
// @Produce json
// @Success 200 {object} Wallet
// @Router /api/v1/wallet [get]
// @Security BearerAuth
func (h *Handler) GetWallet(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	wallet, err := h.service.GetWallet(c.Request().Context(), userID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Wallet retrieved", wallet)
}

func (h *Handler) GetHistory(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	txs, total, err := h.service.GetHistory(c.Request().Context(), userID, page, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	if page < 1 {
		page = 1
	}
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

	return response.SuccessWithMeta(c, http.StatusOK, "Transactions retrieved", txs, meta)
}

// RequestRedeem godoc
// @Summary Request Wallet Redeem
// @Description Submit a cash withdrawal / e-wallet redeem request
// @Tags wallet
// @Accept json
// @Produce json
// @Success 201 {object} RedeemRequest
// @Param request body CreateRedeemDTO true "Redeem Payload"
// @Router /api/v1/wallet/redeem [post]
// @Security BearerAuth
func (h *Handler) RequestRedeem(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var dto CreateRedeemDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	redeem, err := h.service.RequestRedeem(c.Request().Context(), userID, &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Redeem requested successfully", redeem)
}

func (h *Handler) ListRedeems(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	status := c.QueryParam("status")

	items, total, err := h.service.ListRedeems(c.Request().Context(), page, limit, status)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	if page < 1 {
		page = 1
	}
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

	return response.SuccessWithMeta(c, http.StatusOK, "Redeem requests retrieved", items, meta)
}

func (h *Handler) ProcessRedeem(c echo.Context) error {
	idParam := c.Param("id")
	redeemID, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid redeem ID", nil)
	}

	var dto ProcessRedeemDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	res, err := h.service.ProcessRedeem(c.Request().Context(), redeemID, &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, "Redeem request processed", res)
}
