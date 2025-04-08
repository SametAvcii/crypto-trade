package server

import (
	"fmt"
	"net"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func LaunchConsumerServer(appc config.Consumer) {
	app := gin.New()

	app.Use(gin.Logger())
	app.Use(gin.Recovery())
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
	app.Use(otelgin.Middleware(appc.Name))
	fmt.Println("Server is running on port " + appc.Port)
	app.Run(net.JoinHostPort(appc.Host, appc.Port))

}
