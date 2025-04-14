package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Технические метрики
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество HTTP-запросов",
		},
		[]string{"method", "path"},
	)

	ResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Время ответа HTTP-запросов",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Бизнесовые метрики
	CreatedPvz = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pvz_created_total",
			Help: "Количество созданных ПВЗ",
		},
	)

	CreatedReceptions = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "receptions_created_total",
			Help: "Количество созданных приёмок заказов",
		},
	)

	ProductsAdded = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "products_added_total",
			Help: "Количество добавленных товаров",
		},
	)
)

func Init() {
	prometheus.MustRegister(RequestCount, ResponseDuration, CreatedPvz, CreatedReceptions, ProductsAdded)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9000", nil) // Порт 9000
}
