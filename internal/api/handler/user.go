package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pvz/internal/api/mapper"
	"pvz/internal/api/response"
)

func (h *Handler) DummyLogin(c *gin.Context) {
	var req response.DummyLoginPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("Invalid input data for dummy login", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
		return
	}

	token, err := h.service.DummyLogin(c, req.Role)
	if err != nil {
		h.logger.Warnw("Dummy login failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to generate token"})
		return
	}

	h.logger.Infow("Dummy login successful", "role", req.Role)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Register(c *gin.Context) {
	var req response.RegisterPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("Invalid input data for registration", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
		return
	}

	user := mapper.ToUser(req)

	createdUser, err := h.service.CreateUser(c, user)
	if err != nil {
		h.logger.Errorw("User registration failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	resp := mapper.ToRegisterResponse(createdUser)
	h.logger.Infow("User registered successfully", "userID", createdUser.Id, "email", createdUser.Email)

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Login(c *gin.Context) {
	var req response.LoginPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("Invalid input data for login", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
		return
	}

	token, err := h.service.LoginUser(c, req.Email, req.Password)
	if err != nil {
		h.logger.Warnw("Login failed", "email", req.Email, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	h.logger.Infow("Login successful", "email", req.Email)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
