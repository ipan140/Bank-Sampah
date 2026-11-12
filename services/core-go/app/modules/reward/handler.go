package reward

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

func (h *Handler) CreateReward(c echo.Context) error {
	var dto CreateRewardDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	item, err := h.service.CreateReward(c.Request().Context(), &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Reward created successfully", item)
}

func (h *Handler) ListRewards(c echo.Context) error {
	activeOnly := c.QueryParam("all") != "true"
	items, err := h.service.ListRewards(c.Request().Context(), activeOnly)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Rewards retrieved", items)
}

func (h *Handler) ClaimReward(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	idParam := c.Param("id")
	rewardID, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid reward ID", nil)
	}

	claim, err := h.service.ClaimReward(c.Request().Context(), userID, rewardID)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, "Reward claimed successfully", claim)
}

func (h *Handler) MyClaims(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	claims, err := h.service.GetUserClaims(c.Request().Context(), userID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Claim history retrieved", claims)
}
