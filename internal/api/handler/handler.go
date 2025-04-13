package handler

import (
	"github.com/gin-gonic/gin"
	"pvz/internal/logger"
	"pvz/internal/middleware/jwt"
	"pvz/internal/service"
)

type Handler struct {
	service *service.Service
	logger  logger.Logger
}

func NewHandler(services *service.Service, log logger.Logger) *Handler {
	return &Handler{
		service: services,
		logger:  log,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/dummyLogin", h.DummyLogin)
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
	router.POST("/pvz", jwt.AuthMiddleware("moderator"), h.createPvz)
	router.POST("/receptions", jwt.AuthMiddleware("employee"), h.createReception)
	router.POST("/products", jwt.AuthMiddleware("employee"), h.addProduct)
	router.DELETE("/pvz/:pvzId/delete_last_product", jwt.AuthMiddleware("employee"), h.deleteLastProduct)
	router.PATCH("/pvz/:pvzId/close_last_reception", jwt.AuthMiddleware("employee"), h.closeReception)
	router.GET("/pvz", jwt.AuthMiddleware("moderator", "employee"), h.getPvz)

	return router
}
