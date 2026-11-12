package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c echo.Context, statusCode int, message string, data interface{}, meta interface{}) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c echo.Context, statusCode int, message string, errs interface{}) error {
	return c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

func BadRequest(c echo.Context, message string, errs interface{}) error {
	return Error(c, http.StatusBadRequest, message, errs)
}

func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c echo.Context, message string) error {
	return Error(c, http.StatusForbidden, message, nil)
}

func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message, nil)
}

func InternalServerError(c echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, message, nil)
}
