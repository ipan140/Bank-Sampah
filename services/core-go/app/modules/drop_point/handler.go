package drop_point

import (
	"net/http"

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

// ListDropPoints godoc
// @Summary List Drop-Off Points
// @Description Find nearby waste bank drop points and smart trash boxes
// @Tags drop_point
// @Produce json
// @Success 200 {array} DropPoint
// @Router /api/v1/drop-points [get]
func (h *Handler) ListDropPoints(c echo.Context) error {
	city := c.QueryParam("city")
	points, err := h.service.ListDropPoints(c.Request().Context(), city)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}
	return response.Success(c, http.StatusOK, "Drop points retrieved", points)
}

// CreateDropPoint godoc
// @Summary Create Drop Point
// @Description Admin creates a new drop point location
// @Tags drop_point
// @Accept json
// @Produce json
// @Success 201 {object} DropPoint
// @Param request body CreateDropPointDTO true "Drop Point Payload"
// @Router /api/v1/drop-points [post]
// @Security BearerAuth
func (h *Handler) CreateDropPoint(c echo.Context) error {
	var dto CreateDropPointDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	point, err := h.service.CreateDropPoint(c.Request().Context(), &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Drop point created", point)
}

// GetDropPoint godoc
// @Summary Get Drop Point Details
// @Description View specific drop point address and operating hours
// @Tags drop_point
// @Produce json
// @Success 200 {object} DropPoint
// @Router /api/v1/drop-points/{id} [get]
func (h *Handler) GetDropPoint(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid drop point ID", nil)
	}

	point, err := h.service.GetDropPointByID(c.Request().Context(), uid)
	if err != nil {
		return response.NotFound(c, "Drop point not found")
	}

	return response.Success(c, http.StatusOK, "Drop point details retrieved", point)
}

// GetSchedules godoc
// @Summary Get Mobile Truck Schedules
// @Description View weekly truck arrival schedules for RT/RW posko drop-off
// @Tags drop_point
// @Produce json
// @Success 200 {array} MobilePickupSchedule
// @Router /api/v1/drop-points/schedules [get]
func (h *Handler) GetSchedules(c echo.Context) error {
	schedules, err := h.service.GetSchedules(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}
	return response.Success(c, http.StatusOK, "Mobile pickup schedules retrieved", schedules)
}
