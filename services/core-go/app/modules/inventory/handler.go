package inventory

import (
	"net/http"
	"strconv"

	"bank-sampah/core-go/common/response"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetStocks godoc
// @Summary Get Warehouse Waste Stock
// @Description View live physical inventory stocks per category in the warehouse
// @Tags inventory
// @Produce json
// @Success 200 {array} WasteStock
// @Router /api/v1/inventory/stocks [get]
// @Security BearerAuth
func (h *Handler) GetStocks(c echo.Context) error {
	stocks, err := h.service.GetStocks(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}
	return response.Success(c, http.StatusOK, "Warehouse stocks retrieved", stocks)
}

// CreateB2BSale godoc
// @Summary Create B2B Waste Sale to Offtaker Factory
// @Description Sell processed waste batches to recycler factories and record revenue/gross profit margin
// @Tags inventory
// @Accept json
// @Produce json
// @Success 201 {object} B2BSaleOrder
// @Param request body CreateB2BSaleDTO true "Sale Payload"
// @Router /api/v1/inventory/sales [post]
// @Security BearerAuth
func (h *Handler) CreateB2BSale(c echo.Context) error {
	var dto CreateB2BSaleDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	sale, err := h.service.CreateB2BSale(c.Request().Context(), &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "B2B sale completed and stock deducted", sale)
}

// ListB2BSales godoc
// @Summary List B2B Offtaker Sales History
// @Description View historical waste sales to recycling factories
// @Tags inventory
// @Produce json
// @Success 200 {array} B2BSaleOrder
// @Router /api/v1/inventory/sales [get]
// @Security BearerAuth
func (h *Handler) ListB2BSales(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	sales, total, err := h.service.ListB2BSales(c.Request().Context(), page, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	if page < 1 {
		page = 1
	}
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

	return response.SuccessWithMeta(c, http.StatusOK, "B2B sales retrieved", sales, meta)
}

// GetMutations godoc
// @Summary Get Warehouse Stock Mutations
// @Description View stock movements (in from scale, out to factory sale)
// @Tags inventory
// @Produce json
// @Success 200 {array} StockMutation
// @Router /api/v1/inventory/mutations [get]
// @Security BearerAuth
func (h *Handler) GetMutations(c echo.Context) error {
	category := c.QueryParam("category")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	mutations, err := h.service.GetMutations(c.Request().Context(), category, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Stock mutations retrieved", mutations)
}
