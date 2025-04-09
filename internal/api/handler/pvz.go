package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pvz/internal/api/response"
	"pvz/internal/middleware/mapper"
)

func (h *Handler) createPvz(c *gin.Context) {
	var req response.PvzRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	pvz := mapper.ToPvz(req)

	createdPvz, err := h.service.CreatePvz(c.Request.Context(), pvz)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	PvzResponse := mapper.ToPvzResponse(createdPvz)

	c.JSON(http.StatusCreated, PvzResponse)
}
