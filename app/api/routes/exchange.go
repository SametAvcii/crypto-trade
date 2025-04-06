package routes

import (
	"fmt"
	"net/http"

	ctlog "github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/domains/exchange"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/gin-gonic/gin"
)

func ExchangeRoutes(r *gin.RouterGroup, s exchange.Service) {
	r.POST("/", AddExchange(s))

}

// @Summary Add Symbol
// @Description  Add Symbol Route
// @Tags Symbol Endpoints
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param			payload	body	dtos.AddSymbolReq	true	"Add Symbol Request"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /x [GET]
func AddExchange(s exchange.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("Add Exchange")
		var req dtos.AddExchangeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error":  err.Error(),
				"status": 400,
			})
			return
		}

		res, err := s.AddExchange(req)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Add Exchange Error",
				Message: "Add Exchange err: " + err.Error(),
				Entity:  "symbol",
				Type:    "error",
			})
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":  err.Error(),
				"status": http.StatusBadRequest,
			})
			return
		}

		ctlog.CreateLog(&entities.Log{
			Title:   "Add Exchange",
			Message: "Add exchange success: " + res.Name,
			Entity:  "symbol",
			Type:    "success",
		})

		c.JSON(201, gin.H{
			"data":   res,
			"status": 201,
		})
	}
}
