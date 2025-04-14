package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"pvz/internal/api/mapper"
	"pvz/internal/api/response"
)

func (h *Handler) CreateReception(c *gin.Context) {
	var req response.ReceptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("Invalid input data for reception creation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
		return
	}

	if _, err := uuid.Parse(req.PvzId); err != nil {
		h.logger.Warnw("Invalid input data for reception creation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"error":   "PvzId is not a valid UUID",
		})
		return
	}

	reception := mapper.ToReception(req)

	createdReception, err := h.service.CreateReception(c.Request.Context(), reception.PvzId)
	if err != nil {
		h.logger.Errorw("Failed to create reception", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create reception", "error": err.Error()})
		return
	}

	resp := mapper.ToReceptionResponse(createdReception)
	h.logger.Infow("Reception created successfully", "receptionId", createdReception.Id, "pvzId", createdReception.PvzId)

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) CloseReception(c *gin.Context) {
	pvzIdParam := c.Param("pvzId")
	pvzId, err := uuid.Parse(pvzIdParam)
	if err != nil {
		h.logger.Errorw("Invalid PvzId format", "PvzId", pvzIdParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PvzId format"})
		return
	}

	h.logger.Infow("Attempting to close reception", "PvzId", pvzId)

	err = h.service.CloseReception(c, pvzId)
	if err != nil {
		h.logger.Errorw("Failed to close reception", "PvzId", pvzId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to close reception: %v", err)})
		return
	}

	h.logger.Infow("Reception closed successfully", "PvzId", pvzId)

	c.JSON(http.StatusOK, gin.H{"message": "Reception closed successfully"})
}
