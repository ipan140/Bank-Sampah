package user

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	userGroup := g.Group("/users", middleware.JWTMiddleware(cfg))

	// Self profile
	userGroup.PUT("/me", h.UpdateProfile)

	// Admin only
	adminGroup := userGroup.Group("", middleware.RequireRoles(string(auth.RoleAdmin)))
	adminGroup.GET("", h.ListUsers)
	adminGroup.GET("/:id", h.GetUser)
	adminGroup.PATCH("/:id/role", h.AssignRole)
}
