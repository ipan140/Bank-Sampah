package gamification

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

// GetLeaderboard godoc
// @Summary Get Environmental Leaderboard Ranking
// @Description View top ranking RT, RW, schools, or individual citizens by total waste collected
// @Tags gamification
// @Produce json
// @Success 200 {array} LeaderboardEntry
// @Router /api/v1/gamification/leaderboard [get]
func (h *Handler) GetLeaderboard(c echo.Context) error {
	entityType := EntityType(c.QueryParam("type"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	entries, err := h.service.GetLeaderboard(c.Request().Context(), entityType, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Leaderboard retrieved", entries)
}

// GetMyLevel godoc
// @Summary Get My Eco Level and Badges
// @Description View personal recycling level, milestone badges, and progress bar
// @Tags gamification
// @Produce json
// @Success 200 {object} UserLevelProgress
// @Router /api/v1/gamification/my-level [get]
// @Security BearerAuth
func (h *Handler) GetMyLevel(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	level, err := h.service.GetUserLevelAndBadges(c.Request().Context(), userID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Level and badges retrieved", level)
}
