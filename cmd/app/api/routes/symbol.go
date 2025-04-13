package routes

import (
	"net/http"

	ctlog "github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/domains/symbol"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/gin-gonic/gin"
)

func SymbolRoutes(r *gin.RouterGroup, s symbol.Service) {
	r.POST("/", AddSymbol(s))

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
func AddSymbol(s symbol.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dtos.AddSymbolReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error":  err.Error(),
				"status": 400,
			})
			return
		}

		res, err := s.AddSymbol(req)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Add Symbol Error",
				Message: "Add symbol err: " + err.Error(),
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
			Title:   "Add Symbol",
			Message: "Add symbol success: " + res.Symbol,
			Entity:  "symbol",
			Type:    "success",
		})

		c.JSON(201, gin.H{
			"data":   res,
			"status": 201,
		})
	}
}
