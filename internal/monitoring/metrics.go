package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var TTSRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "tts_requests_total",
		Help: "Número total de requisições por provedor",
	},
	[]string{"provider"},
)

var ResponseTime = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "tts_response_time_seconds",
		Help:    "Tempo de resposta das requisições",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"provider"},
)

func InitMetrics() {
	prometheus.MustRegister(TTSRequests)
	prometheus.MustRegister(ResponseTime)
}

func ExposeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9090", nil)
}
