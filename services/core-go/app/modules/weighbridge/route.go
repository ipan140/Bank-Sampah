package weighbridge

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	scaleGroup := g.Group("/weighbridge", middleware.JWTMiddleware(cfg))

	// Admin / Kurir / Operator can create weigh tickets
	scaleGroup.POST("/tickets", h.CreateTicket, middleware.RequireRoles(string(auth.RoleAdmin), string(auth.RoleKurir)))
	scaleGroup.GET("/tickets", h.ListTickets)
	scaleGroup.GET("/tickets/:id", h.GetTicket)
}
