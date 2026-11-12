package tracking

import (
	"net/http"
	"strconv"

	"bank-sampah/core-go/common/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) UpdateLocation(c echo.Context) error {
	courierIDStr := c.Get("user_id").(string)
	courierID, err := uuid.Parse(courierIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid courier ID", nil)
	}

	var dto UpdateLocationDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	if err := h.service.UpdateLocation(c.Request().Context(), courierID, &dto); err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Location updated successfully", nil)
}

func (h *Handler) GetCourierLocation(c echo.Context) error {
	idParam := c.Param("courier_id")
	courierID, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid courier ID", nil)
	}

	loc, err := h.service.GetCourierLiveLocation(c.Request().Context(), courierID)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Courier live location retrieved", loc)
}

func (h *Handler) FindNearbyCouriers(c echo.Context) error {
	latStr := c.QueryParam("lat")
	lngStr := c.QueryParam("lng")
	radiusStr := c.QueryParam("radius")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid latitude", nil)
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid longitude", nil)
	}

	radius, _ := strconv.ParseFloat(radiusStr, 64)
	if radius <= 0 {
		radius = 5.0
	}

	couriers, err := h.service.FindNearbyCouriers(c.Request().Context(), lat, lng, radius)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Nearby couriers retrieved", couriers)
}

func (h *Handler) GetHistory(c echo.Context) error {
	idParam := c.Param("courier_id")
	courierID, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid courier ID", nil)
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	history, err := h.service.GetLocationHistory(c.Request().Context(), courierID, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Location history retrieved", history)
}
