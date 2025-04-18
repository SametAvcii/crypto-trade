package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Depado/ginprom"
	"github.com/SametAvcii/crypto-trade/cmd/app/api/routes"
	"github.com/SametAvcii/crypto-trade/internal/clients/database"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/domains/exchange"
	"github.com/SametAvcii/crypto-trade/pkg/domains/signal"
	"github.com/SametAvcii/crypto-trade/pkg/domains/symbol"
	"github.com/SametAvcii/crypto-trade/pkg/metrics"
	"github.com/SametAvcii/crypto-trade/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var (
	swaggerUser string
	swaggerPass string
)

func init() {
	metrics.Register()
}

func LaunchHttpServer(appc config.App, allows config.Allows) {
	log.Println("Starting HTTP Server...")
	gin.SetMode(gin.ReleaseMode)

	app := gin.New()
	app.Use(middleware.PrometheusMiddleware())

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
	app.Use(gin.Recovery())
	app.Use(otelgin.Middleware(appc.Name))

	pgDB := database.PgClient()

	app.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	app.GET("/readyz", func(c *gin.Context) {
		sqlDB, err := pgDB.DB()
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

	app.Use(cors.New(cors.Config{
		AllowMethods:     allows.Methods,
		AllowHeaders:     allows.Headers,
		AllowOrigins:     allows.Origins,
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	app.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	p := ginprom.New(
		ginprom.Engine(app),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
		ginprom.Ignore("/swagger/*any"),
	)
	app.Use(p.Instrument())

	api := app.Group("/api/v1")
	symbolRoute := api.Group("/symbol")
	symbolRepo := symbol.NewRepo(pgDB)
	symbolService := symbol.NewService(symbolRepo)
	routes.SymbolRoutes(symbolRoute, symbolService)

	exchangeRoute := api.Group("/exchange")
	exchangeRepo := exchange.NewRepo(pgDB)
	exchangeService := exchange.NewService(exchangeRepo)
	routes.ExchangeRoutes(exchangeRoute, exchangeService)

	signalRoute := api.Group("/signal")
	signalRepo := signal.NewRepo(pgDB)
	signalService := signal.NewService(signalRepo)
	routes.SignalRoutes(signalRoute, signalService)

	app.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "docs/index.html")
	})

	if os.Getenv("SWAGGER_USER") != "" {
		swaggerUser = os.Getenv("SWAGGER_USER")
	} else {
		swaggerUser = "crypto-trade-dev"
	}
	if os.Getenv("SWAGGER_PASS") != "" {
		swaggerPass = os.Getenv("SWAGGER_PASS")
	} else {
		swaggerPass = "crypto-trade-dev"
	}

	docs.SwaggerInfo.Host = config.InitConfig().App.BaseUrl
	docs.SwaggerInfo.Version = os.Getenv("APP_VERSION")
	app.GET("/docs/*any",
		gin.BasicAuth(gin.Accounts{
			swaggerUser: swaggerPass,
		}),
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	fmt.Println("Server is running on port " + appc.Port)
	app.Run(net.JoinHostPort(appc.Host, appc.Port))
}
