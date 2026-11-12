package csr

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	csrGroup := g.Group("/csr")

	// Public impact overview
	csrGroup.GET("/impact", h.GetImpactMetrics)
	csrGroup.GET("/programs", h.ListPrograms)
	csrGroup.GET("/programs/:id", h.GetProgram)

	// Admin / Mitra corporate create program
	protected := csrGroup.Group("", middleware.JWTMiddleware(cfg))
	protected.POST("/programs", h.CreateProgram, middleware.RequireRoles(string(auth.RoleAdmin), string(auth.RoleMitra)))
}
