package gamification

import (
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	gameGroup := g.Group("/gamification")

	// Public Leaderboard
	gameGroup.GET("/leaderboard", h.GetLeaderboard)

	// User Level & Badges
	userGroup := gameGroup.Group("", middleware.JWTMiddleware(cfg))
	userGroup.GET("/my-level", h.GetMyLevel)
}
