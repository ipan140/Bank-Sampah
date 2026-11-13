package marketplace

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	marketGroup := g.Group("/marketplace")

	// Public catalog
	marketGroup.GET("/products", h.ListProducts)
	marketGroup.GET("/products/:id", h.GetProduct)

	// User checkout & order history
	userGroup := marketGroup.Group("", middleware.JWTMiddleware(cfg))
	userGroup.POST("/checkout", h.Checkout)
	userGroup.GET("/my-orders", h.GetMyOrders)

	// Admin product management
	adminGroup := userGroup.Group("", middleware.RequireRoles(string(auth.RoleAdmin)))
	adminGroup.POST("/products", h.CreateProduct)
}
