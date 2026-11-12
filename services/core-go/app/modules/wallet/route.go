package wallet

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	walletGroup := g.Group("/wallet", middleware.JWTMiddleware(cfg))

	walletGroup.GET("", h.GetWallet)
	walletGroup.GET("/transactions", h.GetHistory)
	walletGroup.POST("/redeem", h.RequestRedeem)

	// Admin only
	adminGroup := walletGroup.Group("", middleware.RequireRoles(string(auth.RoleAdmin)))
	adminGroup.GET("/redeems", h.ListRedeems)
	adminGroup.PATCH("/redeems/:id", h.ProcessRedeem)
}
