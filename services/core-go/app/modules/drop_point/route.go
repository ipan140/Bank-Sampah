package drop_point

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	dpGroup := g.Group("/drop-points")

	// Public access for map search
	dpGroup.GET("", h.ListDropPoints)
	dpGroup.GET("/schedules", h.GetSchedules)
	dpGroup.GET("/:id", h.GetDropPoint)

	// Admin create
	adminGroup := dpGroup.Group("", middleware.JWTMiddleware(cfg), middleware.RequireRoles(string(auth.RoleAdmin)))
	adminGroup.POST("", h.CreateDropPoint)
}
