package member_card

import (
	"net/http"

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

// GetMyCard godoc
// @Summary Get Digital Member Card & QR Code
// @Description Retrieve the authenticated user's digital membership card and QR code token
// @Tags member_card
// @Produce json
// @Success 200 {object} MemberCard
// @Router /api/v1/member-card/me [get]
// @Security BearerAuth
func (h *Handler) GetMyCard(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	card, err := h.service.GetMyCard(c.Request().Context(), userID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Member card retrieved", card)
}

// ScanMemberQR godoc
// @Summary Scan Nasabah QR Code (For Couriers/Weighbridge Operators)
// @Description Quick lookup of member info and balance by scanning their QR code
// @Tags member_card
// @Accept json
// @Produce json
// @Success 200 {object} MemberInfoResponse
// @Param request body ScanQRRequest true "QR Token"
// @Router /api/v1/member-card/scan [post]
// @Security BearerAuth
func (h *Handler) ScanMemberQR(c echo.Context) error {
	var req ScanQRRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	info, err := h.service.ScanMemberQR(c.Request().Context(), req.QRCodeToken)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Member info retrieved successfully", info)
}
