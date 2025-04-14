package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"pvz/internal/logger"
	"pvz/internal/middleware/jwt"
	"pvz/internal/service"
	"pvz/metrics"
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

	router.POST("/dummyLogin", h.trackMetrics(h.DummyLogin))
	router.POST("/register", h.trackMetrics(h.Register))
	router.POST("/login", h.trackMetrics(h.Login))
	router.POST("/pvz", jwt.AuthMiddleware("moderator"), h.trackMetrics(h.CreatePvz))
	router.POST("/receptions", jwt.AuthMiddleware("employee"), h.trackMetrics(h.CreateReception))
	router.POST("/products", jwt.AuthMiddleware("employee"), h.trackMetrics(h.AddProduct))
	router.DELETE("/pvz/:pvzId/delete_last_product", jwt.AuthMiddleware("employee"), h.trackMetrics(h.DeleteLastProduct))
	router.PATCH("/pvz/:pvzId/close_last_reception", jwt.AuthMiddleware("employee"), h.trackMetrics(h.CloseReception))
	router.GET("/pvz", jwt.AuthMiddleware("moderator", "employee"), h.trackMetrics(h.GetPvz))

	return router
}

func (h *Handler) trackMetrics(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Записываем метрики до выполнения запроса
		start := time.Now()

		// Выполняем хендлер
		handlerFunc(c)

		// Измеряем продолжительность запроса
		duration := time.Since(start).Seconds()

		// Инкрементируем технические метрики
		metrics.RequestCount.WithLabelValues(c.Request.Method, c.FullPath()).Inc()
		metrics.ResponseDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
	}
}
