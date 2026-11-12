package csr

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

func (h *Handler) CreateProgram(c echo.Context) error {
	var dto CreateProgramDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	program, err := h.service.CreateProgram(c.Request().Context(), &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "CSR program created", program)
}

func (h *Handler) ListPrograms(c echo.Context) error {
	programs, err := h.service.ListPrograms(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "CSR programs retrieved", programs)
}

func (h *Handler) GetProgram(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid program ID", nil)
	}

	program, err := h.service.GetProgramByID(c.Request().Context(), uid)
	if err != nil {
		return response.NotFound(c, "CSR program not found")
	}

	return response.Success(c, http.StatusOK, "CSR program retrieved", program)
}

func (h *Handler) GetImpactMetrics(c echo.Context) error {
	impact, err := h.service.GetImpactMetrics(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "ESG impact summary retrieved", impact)
}
