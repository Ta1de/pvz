package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"pvz/internal/api/mapper"
	"pvz/internal/api/response"
)

func (h *Handler) AddProduct(c *gin.Context) {
	var req response.ProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind product request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	pvzId, err := uuid.Parse(req.PvzId)
	if err != nil {
		h.logger.Errorw("Invalid PvzId format", "PvzId", req.PvzId, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PvzId format"})
		return
	}

	product := mapper.ToProduct(req)

	createdProduct, err := h.service.AddProduct(c, pvzId, product.Type)
	if err != nil {
		h.logger.Errorw("Failed to add product", "error", err, "PvzId", pvzId, "type", product.Type)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
		return
	}

	productResponse := mapper.ToProductResponse(createdProduct)

	c.JSON(http.StatusOK, productResponse)
}

func (h *Handler) DeleteLastProduct(c *gin.Context) {
	pvzIdParam := c.Param("pvzId")
	pvzId, err := uuid.Parse(pvzIdParam)
	if err != nil {
		h.logger.Errorw("Invalid PvzId format", "PvzId", pvzIdParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PvzId format"})
		return
	}

	h.logger.Infow("Attempting to delete last product", "PvzId", pvzId)

	err = h.service.DeleteLastProduct(c, pvzId)
	if err != nil {
		h.logger.Errorw("Failed to delete last product", "PvzId", pvzId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete last product: %v", err)})
		return
	}

	h.logger.Infow("Last product deleted successfully", "PvzId", pvzId)

	c.JSON(http.StatusOK, gin.H{"message": "Last product deleted successfully"})
}
