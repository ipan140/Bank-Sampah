package ppob

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

// PayBill godoc
// @Summary Pay Utility Bills with Waste Wallet
// @Description Purchase PLN electricity token, mobile phone credit/data, BPJS, PDAM, or PBB tax using waste bank balance
// @Tags ppob
// @Accept json
// @Produce json
// @Success 201 {object} PPOBTransaction
// @Param request body PayBillDTO true "Bill Payload"
// @Router /api/v1/ppob/pay [post]
// @Security BearerAuth
func (h *Handler) PayBill(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var dto PayBillDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	bill, err := h.service.PayBill(c.Request().Context(), userID, &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Bill payment successful", bill)
}

// GetGoldAccount godoc
// @Summary Get Gold Savings Balance
// @Description View digital gold savings balance (in grams) converted from recycled waste earnings
// @Tags ppob
// @Produce json
// @Success 200 {object} GoldSavingsAccount
// @Router /api/v1/ppob/gold [get]
// @Security BearerAuth
func (h *Handler) GetGoldAccount(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	acc, err := h.service.GetGoldAccount(c.Request().Context(), userID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Gold account retrieved", acc)
}

// ConvertToGold godoc
// @Summary Convert Recycled Waste Balance to Digital Gold
// @Description Invest waste bank wallet balance into digital gold savings
// @Tags ppob
// @Accept json
// @Produce json
// @Success 200 {object} GoldSavingsAccount
// @Param request body ConvertGoldDTO true "Convert Payload"
// @Router /api/v1/ppob/gold/convert [post]
// @Security BearerAuth
func (h *Handler) ConvertToGold(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var dto ConvertGoldDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	acc, err := h.service.ConvertBalanceToGold(c.Request().Context(), userID, dto.AmountIDR)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, "Balance successfully converted to gold", acc)
}

// CreateDonation godoc
// @Summary Donate Waste Earnings to Social Charities
// @Description Sedekah / Donate waste bank balance directly to orphanages or mosques
// @Tags ppob
// @Accept json
// @Produce json
// @Success 201 {object} WasteDonation
// @Param request body CreateDonationDTO true "Donation Payload"
// @Router /api/v1/ppob/donate [post]
// @Security BearerAuth
func (h *Handler) CreateDonation(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var dto CreateDonationDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	donation, err := h.service.CreateDonation(c.Request().Context(), userID, &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Donation completed successfully", donation)
}

// GetPPOBHistory godoc
// @Summary Get PPOB Transaction History
// @Description View history of bill payments
// @Tags ppob
// @Produce json
// @Success 200 {array} PPOBTransaction
// @Router /api/v1/ppob/history [get]
// @Security BearerAuth
func (h *Handler) GetPPOBHistory(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	history, err := h.service.GetPPOBHistory(c.Request().Context(), userID, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "PPOB history retrieved", history)
}
