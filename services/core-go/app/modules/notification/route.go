package notification

import (
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	notifGroup := g.Group("/notifications", middleware.JWTMiddleware(cfg))

	notifGroup.GET("", h.GetMyNotifications)
	notifGroup.PATCH("/:id/read", h.MarkAsRead)
}
