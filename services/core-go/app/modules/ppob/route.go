package ppob

import (
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	ppobGroup := g.Group("/ppob", middleware.JWTMiddleware(cfg))

	ppobGroup.POST("/pay", h.PayBill)
	ppobGroup.GET("/history", h.GetPPOBHistory)
	ppobGroup.GET("/gold", h.GetGoldAccount)
	ppobGroup.POST("/gold/convert", h.ConvertToGold)
	ppobGroup.POST("/donate", h.CreateDonation)
}
