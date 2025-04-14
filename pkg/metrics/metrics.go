package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_requests_total",
			Help: "Total number of HTTP requests processed.",
		},
		[]string{"method", "status"},
	)
	HttpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_error_total",
			Help: "Total number of HTTP errors.",
		},
		[]string{"path", "method"},
	)

	LatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_latency_seconds",
			Help:    "Request latency distributions.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	ResponseSizeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "Response size distributions.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"path"},
	)

	RequestSizeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_size_bytes",
			Help:    "Request size distributions.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"path"},
	)

	StatusCodeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "status_code_total",
			Help: "HTTP response status codes.",
		},
		[]string{"status"},
	)

	RequestMethodCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_method_total",
			Help: "HTTP method counts.",
		},
		[]string{"method"},
	)

	RequestPathCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_path_total",
			Help: "Request path counts.",
		},
		[]string{"path"},
	)
)

func Register() {
	prometheus.MustRegister(
		HttpRequests,
		HttpDuration,
		ErrorCounter,
		LatencyHistogram,
		ResponseSizeHistogram,
		RequestSizeHistogram,
		StatusCodeCounter,
		RequestMethodCounter,
		RequestPathCounter,
	)
}
