package finance_report

import (
	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler, cfg *config.Config) {
	finGroup := g.Group("/finance", middleware.JWTMiddleware(cfg), middleware.RequireRoles(string(auth.RoleAdmin)))

	finGroup.GET("/profit-and-loss", h.GetProfitAndLoss)
	finGroup.POST("/expenses", h.RecordExpense)
	finGroup.GET("/expenses", h.GetExpenses)
}
