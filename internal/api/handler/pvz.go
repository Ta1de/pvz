package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pvz/internal/api/response"
	"pvz/internal/logger"
	"pvz/internal/middleware/mapper"
)

func (h *Handler) createPvz(c *gin.Context) {
	var req response.PvzRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SugaredLogger.Warnw("Invalid PvzRequest", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	logger.SugaredLogger.Infow("Creating new PVZ", "request", req)

	pvz := mapper.ToPvz(req)

	createdPvz, err := h.service.CreatePvz(c.Request.Context(), pvz)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to create PVZ", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	PvzResponse := mapper.ToPvzResponse(createdPvz)

	logger.SugaredLogger.Infow("Successfully created PVZ", "pvz", PvzResponse)
	c.JSON(http.StatusCreated, PvzResponse)
}
