package middleware

import (
	"bank-sampah/core-go/common/response"

	"github.com/labstack/echo/v4"
)

func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleVal := c.Get("role")
			if roleVal == nil {
				return response.Forbidden(c, "Forbidden: No role assigned")
			}

			userRole, ok := roleVal.(string)
			if !ok {
				return response.Forbidden(c, "Forbidden: Invalid role")
			}

			for _, allowed := range allowedRoles {
				if userRole == allowed {
					return next(c)
				}
			}

			return response.Forbidden(c, "Forbidden: Insufficient privileges")
		}
	}
}
