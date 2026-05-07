package http

import (
	"strconv"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
	"sunchem-backend/internal/modules/analytics/domain"
	"sunchem-backend/internal/modules/analytics/usecase"

	"github.com/gin-gonic/gin"
)

type TrafficHandler struct {
	uc *usecase.TrafficUseCase
}

func NewTrafficHandler(uc *usecase.TrafficUseCase) *TrafficHandler {
	return &TrafficHandler{uc: uc}
}

func (h *TrafficHandler) Track(c *gin.Context) {
	var event domain.TrafficEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if appErr := h.uc.Track(&event); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func (h *TrafficHandler) Summary(c *gin.Context) {
	summary, appErr := h.uc.Summary()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, summary)
}

func (h *TrafficHandler) Detail(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	events, appErr := h.uc.Detail(page, limit)
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, events)
}

func (h *TrafficHandler) PageViews(c *gin.Context) {
	counts, appErr := h.uc.PageViews()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, counts)
}
