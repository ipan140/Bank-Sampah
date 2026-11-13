package compliance

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	compGroup := g.Group("/compliance")

	// Public certificate verification
	compGroup.GET("/certificates/verify/:hash", h.VerifyCertificate)

	// Admin / Corporate protected
	protected := compGroup.Group("", middleware.JWTMiddleware(cfg))
	protected.GET("/certificates", h.ListCertificates)
	protected.GET("/subscriptions", h.ListSubscriptions)

	admin := protected.Group("", middleware.RequireRoles(string(auth.RoleAdmin), string(auth.RoleMitra)))
	admin.POST("/certificates", h.IssueCertificate)
	admin.POST("/subscriptions", h.CreateSubscription)
}
