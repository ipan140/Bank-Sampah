package notification

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

// GetMyNotifications godoc
// @Summary Get My Notifications
// @Description View user in-app notifications
// @Tags notification
// @Produce json
// @Success 200 {array} Notification
// @Router /api/v1/notifications [get]
// @Security BearerAuth
func (h *Handler) GetMyNotifications(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	notifs, err := h.service.GetUserNotifications(c.Request().Context(), userID, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Notifications retrieved", notifs)
}

// MarkAsRead godoc
// @Summary Mark Notification as Read
// @Description Update notification status to READ
// @Tags notification
// @Produce json
// @Success 200 {string} string "OK"
// @Router /api/v1/notifications/{id}/read [patch]
// @Security BearerAuth
func (h *Handler) MarkAsRead(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	idParam := c.Param("id")
	notifID, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid notification ID", nil)
	}

	if err := h.service.MarkAsRead(c.Request().Context(), notifID, userID); err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Notification marked as read", nil)
}
