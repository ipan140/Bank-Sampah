package reward

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	rewardGroup := g.Group("/rewards")

	// Public catalog
	rewardGroup.GET("", h.ListRewards)

	// Protected user operations
	protected := rewardGroup.Group("", middleware.JWTMiddleware(cfg))
	protected.POST("/:id/claim", h.ClaimReward)
	protected.GET("/my-claims", h.MyClaims)

	// Admin reward management
	admin := protected.Group("", middleware.RequireRoles(string(auth.RoleAdmin)))
	admin.POST("", h.CreateReward)
}
