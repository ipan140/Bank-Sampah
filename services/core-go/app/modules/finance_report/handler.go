package finance_report

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

// GetProfitAndLoss godoc
// @Summary Get Profit and Loss Statement
// @Description View bank sampah financial statement: B2B revenue, COGS to nasabah, gross margin, operational expenses, and net profit
// @Tags finance_report
// @Produce json
// @Success 200 {object} ProfitAndLossStatement
// @Router /api/v1/finance/profit-and-loss [get]
// @Security BearerAuth
func (h *Handler) GetProfitAndLoss(c echo.Context) error {
	pnl, err := h.service.GetProfitAndLoss(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}
	return response.Success(c, http.StatusOK, "Profit and Loss statement retrieved", pnl)
}

// RecordExpense godoc
// @Summary Record Operational Expense
// @Description Log operating cash expenses such as fuel, truck repair, packing ropes, sacks, or labor
// @Tags finance_report
// @Accept json
// @Produce json
// @Success 201 {object} OperationalExpense
// @Param request body CreateExpenseDTO true "Expense Payload"
// @Router /api/v1/finance/expenses [post]
// @Security BearerAuth
func (h *Handler) RecordExpense(c echo.Context) error {
	var dto CreateExpenseDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	exp, err := h.service.RecordExpense(c.Request().Context(), &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Expense recorded successfully", exp)
}

// GetExpenses godoc
// @Summary Get Operational Expenses
// @Description View recent operational expenses
// @Tags finance_report
// @Produce json
// @Success 200 {array} OperationalExpense
// @Router /api/v1/finance/expenses [get]
// @Security BearerAuth
func (h *Handler) GetExpenses(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	expenses, err := h.service.GetExpenses(c.Request().Context(), limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}
	return response.Success(c, http.StatusOK, "Expenses retrieved", expenses)
}
