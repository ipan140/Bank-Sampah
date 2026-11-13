package member_card

import (
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	cardGroup := g.Group("/member-card", middleware.JWTMiddleware(cfg))

	cardGroup.GET("/me", h.GetMyCard)
	cardGroup.POST("/scan", h.ScanMemberQR)
}
