package tracking

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	trackingGroup := g.Group("/tracking", middleware.JWTMiddleware(cfg))

	// Courier pings location
	trackingGroup.POST("/ping", h.UpdateLocation, middleware.RequireRoles(string(auth.RoleKurir)))

	// Query live tracking & nearby
	trackingGroup.GET("/courier/:courier_id", h.GetCourierLocation)
	trackingGroup.GET("/nearby", h.FindNearbyCouriers)
	trackingGroup.GET("/courier/:courier_id/history", h.GetHistory, middleware.RequireRoles(string(auth.RoleAdmin)))
}
