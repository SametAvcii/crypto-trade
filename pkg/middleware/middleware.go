package middleware

import (
	"strconv"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/metrics"
	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method
		statusCode := c.Writer.Status()
		responseSize := float64(c.Writer.Size())
		requestSize := float64(c.Request.ContentLength)

		metrics.HttpRequests.WithLabelValues(method, strconv.Itoa(statusCode)).Inc()
		metrics.HttpDuration.WithLabelValues(path, method).Observe(duration)
		metrics.LatencyHistogram.WithLabelValues(path).Observe(duration)
		metrics.ResponseSizeHistogram.WithLabelValues(path).Observe(responseSize)
		metrics.RequestSizeHistogram.WithLabelValues(path).Observe(requestSize)
		metrics.StatusCodeCounter.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		metrics.RequestMethodCounter.WithLabelValues(method).Inc()
		metrics.RequestPathCounter.WithLabelValues(path).Inc()

		if statusCode >= 400 {
			metrics.ErrorCounter.WithLabelValues(path, method).Inc()
		}
	}
}

//maybe write jwt middleware for routes
 