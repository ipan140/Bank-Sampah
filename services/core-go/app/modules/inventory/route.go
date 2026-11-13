package inventory

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	invGroup := g.Group("/inventory", middleware.JWTMiddleware(cfg))

	// Protected for Admin / Operator
	adminGroup := invGroup.Group("", middleware.RequireRoles(string(auth.RoleAdmin)))
	adminGroup.GET("/stocks", h.GetStocks)
	adminGroup.POST("/sales", h.CreateB2BSale)
	adminGroup.GET("/sales", h.ListB2BSales)
	adminGroup.GET("/mutations", h.GetMutations)
}
