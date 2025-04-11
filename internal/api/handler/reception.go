package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"pvz/internal/api/response"
	"pvz/internal/logger"
	"pvz/internal/middleware/mapper"
)

func (h *Handler) createReception(c *gin.Context) {
	var req response.ReceptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SugaredLogger.Warnw("Invalid input data for reception creation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
		return
	}

	reception := mapper.ToReception(req)

	createdReception, err := h.service.CreateReception(c.Request.Context(), reception.PvzId)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to create reception", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create reception", "error": err.Error()})
		return
	}

	resp := mapper.ToReceptionResponse(createdReception)
	logger.SugaredLogger.Infow("Reception created successfully", "receptionId", createdReception.Id, "pvzId", createdReception.PvzId)

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) closeReception(c *gin.Context) {
	pvzIdParam := c.Param("pvzId")
	pvzId, err := uuid.Parse(pvzIdParam)
	if err != nil {
		logger.SugaredLogger.Errorw("Invalid PvzId format", "PvzId", pvzIdParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PvzId format"})
		return
	}

	logger.SugaredLogger.Infow("Attempting to close reception", "PvzId", pvzId)

	err = h.service.CloseReception(c, pvzId)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to close reception", "PvzId", pvzId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to close reception: %v", err)})
		return
	}

	logger.SugaredLogger.Infow("Reception closed successfully", "PvzId", pvzId)

	c.JSON(http.StatusOK, gin.H{"message": "Reception closed successfully"})
}
