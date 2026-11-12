package auth

import (
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	authGroup := g.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)

	// Protected routes
	protected := authGroup.Group("", middleware.JWTMiddleware(cfg))
	protected.POST("/logout", h.Logout)
	protected.GET("/profile", h.Profile)
}
