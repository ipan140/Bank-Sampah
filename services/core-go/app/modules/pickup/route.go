package pickup

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	pickupGroup := g.Group("/pickups", middleware.JWTMiddleware(cfg))

	pickupGroup.POST("", h.CreatePickup)
	pickupGroup.GET("", h.ListPickups)
	pickupGroup.GET("/:id", h.GetPickup)

	// Admin / Kurir status updates
	pickupGroup.PATCH("/:id/status", h.UpdateStatus, middleware.RequireRoles(string(auth.RoleAdmin), string(auth.RoleKurir)))

	// Admin courier assignment
	pickupGroup.PATCH("/:id/assign", h.AssignCourier, middleware.RequireRoles(string(auth.RoleAdmin)))
}
