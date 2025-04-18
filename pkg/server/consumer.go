package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/SametAvcii/crypto-trade/internal/clients/database"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func LaunchConsumerServer(appc config.Consumer) {
	// Gin app'ini başlatıyoruz
	app := gin.New()

	// Logger ve recovery middlewares
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	// Custom log format
	app.Use(gin.LoggerWithFormatter(func(log gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] - %s \"%s %s %s %d %s\"\n",
			log.TimeStamp.Format("2006-01-02 15:04:05"),
			log.ClientIP,
			log.Method,
			log.Path,
			log.Request.Proto,
			log.StatusCode,
			log.Latency,
		)
	}))

	// OpenTelemetry middleware
	app.Use(otelgin.Middleware(appc.Name))

	// Prometheus metrics endpoint
	app.GET("/metrics", gin.WrapH(promhttp.Handler()))

	app.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	db := database.PgClient()

	app.GET("/readyz", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.String(http.StatusInternalServerError, "Database not ready")
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.String(http.StatusInternalServerError, "Database not reachable")
			return
		}
		c.String(http.StatusOK, "READY")
	})

	// Log mesajı
	fmt.Println("Server is running on port " + appc.Port)

	// Sunucuyu çalıştırıyoruz
	app.Run(net.JoinHostPort(appc.Host, appc.Port))
}
