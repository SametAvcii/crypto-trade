package routes

import (
	"net/http"

	ctlog "github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/domains/signal"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/gin-gonic/gin"
)

func SignalRoutes(r *gin.RouterGroup, s signal.Service) {
	r.POST("/interval", AddSignalInterval(s))

}

// @Summary Add Signal Interval
// @Description  Add Signal Interval Route
// @Tags Signal Endpoints
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param			payload	body	dtos.AddSignalIntervalReq	true	"Add Symbol Request"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /signal/interval [POST]
func AddSignalInterval(s signal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dtos.AddSignalIntervalReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error":  err.Error(),
				"status": 400,
			})
			return
		}

		res, err := s.AddSignalIntervals(req)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Add Signal Interval Error",
				Message: "Add Signal Interval err: " + err.Error(),
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
			Title:   "Add Signal Interval",
			Message: "Add Signal Interval success: " + res.Symbol,
			Entity:  "symbol",
			Type:    "success",
		})

		c.JSON(201, gin.H{
			"data":   res,
			"status": 201,
		})
	}
}
