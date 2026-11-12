package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func AuditLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			req := c.Request()
			res := c.Response()

			userID := c.Get("user_id")
			if userID == nil {
				userID = "anonymous"
			}

			log.Printf("[AUDIT] %s | %3d | %10v | %-15s | %-7s %s | user: %v",
				time.Now().Format("2006-01-02 15:04:05"),
				res.Status,
				latency,
				c.RealIP(),
				req.Method,
				req.RequestURI,
				userID,
			)

			return err
		}
	}
}
