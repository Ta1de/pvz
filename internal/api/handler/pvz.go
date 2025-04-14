package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pvz/internal/api/mapper"
	"pvz/internal/api/response"
	"pvz/internal/logger"
)

func (h *Handler) CreatePvz(c *gin.Context) {
	var req response.PvzRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("Invalid PvzRequest", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	h.logger.Infow("Creating new PVZ", "request", req)

	pvz := mapper.ToPvz(req)

	createdPvz, err := h.service.CreatePvz(c.Request.Context(), pvz)
	if err != nil {
		h.logger.Errorw("Failed to create PVZ", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	PvzResponse := mapper.ToPvzResponse(createdPvz)

	h.logger.Infow("Successfully created PVZ", "pvz", PvzResponse)
	c.JSON(http.StatusCreated, PvzResponse)
}

func (h *Handler) GetPvz(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	h.logger.Infow("Received request for Pvz list",
		"limit", limitStr, "offset", offsetStr, "startDate", startDateStr, "endDate", endDateStr)

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.logger.Warnw("Invalid limit", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.logger.Warnw("Invalid offset", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	var startDate, endDate *time.Time

	if startDateStr != "" {
		t, err := ParseFlexibleTime(startDateStr)
		if err != nil {
			h.logger.Warnw("Invalid startDate", "startDate", startDateStr, "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate"})
			return
		}
		startDate = t
	}

	if endDateStr != "" {
		t, err := ParseFlexibleTime(endDateStr)
		if err != nil {
			h.logger.Warnw("Invalid endDate", "endDate", endDateStr, "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate"})
			return
		}
		endDate = t
	}

	result, err := h.service.GetPvzList(c.Request.Context(), limit, offset, startDate, endDate)
	if err != nil {
		h.logger.Errorw("Failed to get Pvz list", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infow("Successfully retrieved Pvz list", "count", len(result))
	c.JSON(http.StatusOK, result)
}

func ParseFlexibleTime(str string) (*time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05.99999",
		"2006-01-02 15:04:05.9999",
		"2006-01-02 15:04:05",
	}
	for _, layout := range formats {
		if t, err := time.Parse(layout, str); err == nil {
			return &t, nil
		}
	}
	logger.Log.Warnw("Failed to parse time", "input", str)
	return nil, fmt.Errorf("cannot parse time: %s", str)
}
