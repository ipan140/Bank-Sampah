package weighbridge

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

func (h *Handler) CreateTicket(c echo.Context) error {
	operatorIDStr := c.Get("user_id").(string)
	operatorID, err := uuid.Parse(operatorIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid operator ID", nil)
	}

	var dto CreateTicketDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	ticket, err := h.service.CreateTicket(c.Request().Context(), operatorID, &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Scale ticket created and balance settled", ticket)
}

func (h *Handler) ListTickets(c echo.Context) error {
	var filter ScaleTicketFilter
	if err := c.Bind(&filter); err != nil {
		return response.BadRequest(c, "Invalid filter params", err.Error())
	}

	tickets, total, err := h.service.ListTickets(c.Request().Context(), filter)
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

	return response.SuccessWithMeta(c, http.StatusOK, "Tickets retrieved", tickets, meta)
}

func (h *Handler) GetTicket(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid ticket ID", nil)
	}

	ticket, err := h.service.GetTicketByID(c.Request().Context(), uid)
	if err != nil {
		return response.NotFound(c, "Ticket not found")
	}

	return response.Success(c, http.StatusOK, "Ticket retrieved", ticket)
}
