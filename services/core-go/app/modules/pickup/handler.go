package pickup

import (
	"net/http"

	"bank-sampah/core-go/app/modules/auth"
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

// CreatePickup godoc
// @Summary Schedule Waste Pickup
// @Description Schedule a new waste pickup request for couriers
// @Tags pickup
// @Accept json
// @Produce json
// @Success 201 {object} PickupRequest
// @Param request body CreatePickupDTO true "Pickup Payload"
// @Router /api/v1/pickups [post]
// @Security BearerAuth
func (h *Handler) CreatePickup(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	nasabahID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var dto CreatePickupDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	res, err := h.service.CreatePickup(c.Request().Context(), nasabahID, &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Pickup scheduled successfully", res)
}

func (h *Handler) ListPickups(c echo.Context) error {
	var filter PickupFilter
	if err := c.Bind(&filter); err != nil {
		return response.BadRequest(c, "Invalid filter parameters", err.Error())
	}

	role := c.Get("role").(string)
	userIDStr := c.Get("user_id").(string)

	// Scoping: Nasabah only sees own pickups, Kurir sees assigned
	if role == string(auth.RoleNasabah) {
		filter.NasabahID = userIDStr
	} else if role == string(auth.RoleKurir) {
		filter.CourierID = userIDStr
	}

	items, total, err := h.service.ListPickups(c.Request().Context(), filter)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	meta := response.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return response.SuccessWithMeta(c, http.StatusOK, "Pickups retrieved", items, meta)
}

func (h *Handler) GetPickup(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid pickup ID", nil)
	}

	item, err := h.service.GetPickupByID(c.Request().Context(), uid)
	if err != nil {
		return response.NotFound(c, "Pickup not found")
	}

	return response.Success(c, http.StatusOK, "Pickup details retrieved", item)
}

func (h *Handler) AssignCourier(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid pickup ID", nil)
	}

	var dto AssignCourierDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	item, err := h.service.AssignCourier(c.Request().Context(), uid, dto.CourierID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Courier assigned successfully", item)
}

func (h *Handler) UpdateStatus(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid pickup ID", nil)
	}

	var dto UpdateStatusDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	item, err := h.service.UpdateStatus(c.Request().Context(), uid, dto.Status)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Pickup status updated", item)
}
