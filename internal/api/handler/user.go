package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pvz/internal/api/response"
	"pvz/internal/logger"
	"pvz/internal/middleware/mapper"
)

func (h *Handler) dummyLogin(c *gin.Context) {
	var req response.DummyLoginPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SugaredLogger.Warnw("Invalid input data for dummy login", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
		return
	}

	token, err := h.service.DummyLogin(c, req.Role)
	if err != nil {
		logger.SugaredLogger.Warnw("Dummy login failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to generate token"})
		return
	}

	logger.SugaredLogger.Infow("Dummy login successful", "role", req.Role)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) register(c *gin.Context) {
	var req response.RegisterPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SugaredLogger.Warnw("Invalid input data for registration", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
		return
	}

	user := mapper.ToUser(req)

	createdUser, err := h.service.CreateUser(c, user)
	if err != nil {
		logger.SugaredLogger.Errorw("User registration failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	resp := mapper.ToRegisterResponse(createdUser)
	logger.SugaredLogger.Infow("User registered successfully", "userID", createdUser.Id, "email", createdUser.Email)

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) login(c *gin.Context) {
	var req response.LoginPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SugaredLogger.Warnw("Invalid input data for login", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
		return
	}

	token, err := h.service.LoginUser(c, req.Email, req.Password)
	if err != nil {
		logger.SugaredLogger.Warnw("Login failed", "email", req.Email, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	logger.SugaredLogger.Infow("Login successful", "email", req.Email)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
