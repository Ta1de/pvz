package handler

import (
	"github.com/gin-gonic/gin"
	"pvz/internal/middleware/jwt"
	"pvz/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/dummyLogin", h.dummyLogin)
	router.POST("/register", h.register)
	router.POST("/login", h.login)
	router.POST("/pvz", jwt.AuthMiddleware("moderator"), h.createPvz)
	router.POST("/reception", jwt.AuthMiddleware("employee"), h.createReception)

	return router
}
