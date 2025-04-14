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
	r.PUT("/:id", UpdateSymbol(s))
	r.DELETE("/:id", DeleteSymbol(s))
	r.GET("/", GetAllSymbols(s))
	r.GET("/:id", GetSymbolByID(s))

}

// @Summary Add Symbol
// @Description  Add Symbol Route
// @Tags Symbol Endpoints
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param			payload	body	dtos.AddSymbolReq	true	"Add Symbol Request"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router / [POST]
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

		res, err := s.AddSymbol(c, req)
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

// @Summary Update Symbol
// @Description Update Symbol Route
// @Tags Symbol Endpoints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Symbol ID"
// @Param payload body dtos.UpdateSymbolReq true "Update Symbol Request"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /symbol/{id} [PUT]
func UpdateSymbol(s symbol.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dtos.UpdateSymbolReq
		id := c.Param("id")

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error":  err.Error(),
				"status": 400,
			})
			return
		}
		req.ID = id

		res, err := s.UpdateSymbol(c, req)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Update Symbol Error",
				Message: "Update symbol err: " + err.Error(),
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
			Title:   "Update Symbol",
			Message: "Update symbol success: " + res.Symbol,
			Entity:  "symbol",
			Type:    "success",
		})

		c.JSON(200, gin.H{
			"data":   res,
			"status": 200,
		})
	}
}

// @Summary Delete Symbol
// @Description Delete Symbol Route
// @Tags Symbol Endpoints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Symbol ID"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /symbol/{id} [DELETE]
func DeleteSymbol(s symbol.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := s.DeleteSymbol(c, id)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Delete Symbol Error",
				Message: "Delete symbol err: " + err.Error(),
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
			Title:   "Delete Symbol",
			Message: "Delete symbol success",
			Entity:  "symbol",
			Type:    "success",
		})

		c.JSON(200, gin.H{
			"message": "Symbol deleted successfully",
			"status":  200,
		})
	}
}

// @Summary Get All Symbols
// @Description Get All Symbols Route
// @Tags Symbol Endpoints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dtos.GetSymbolRes
// @Failure 400 {object} map[string]any
// @Router /symbol [GET]
func GetAllSymbols(s symbol.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		symbols, err := s.GetAllSymbols(c)
		if err != nil {

			ctlog.CreateLog(&entities.Log{
				Title:   "Get All Symbols Error",
				Message: "Get all symbols err: " + err.Error(),
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
			Title:   "Get All Symbols",
			Message: "Get all symbols success",
			Entity:  "symbol",
			Type:    "success",
		})

		c.JSON(200, gin.H{
			"data":   symbols,
			"status": 200,
		})
	}
}

// @Summary Get Symbol by ID
// @Description Get Symbol by ID Route
// @Tags Symbol Endpoints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Symbol ID"
// @Success 200 {object} dtos.GetSymbolRes
// @Failure 400 {object} map[string]any
// @Router /symbol/{id} [GET]
func GetSymbolByID(s symbol.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		symbol, err := s.GetSymbol(c, id)
		if err != nil {

			ctlog.CreateLog(&entities.Log{
				Title:   "Get Symbol Error",
				Message: "Get symbol err: " + err.Error(),
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
			Title:   "Get Symbol",
			Message: "Get symbol success: " + symbol.Symbol,
			Entity:  "symbol",
			Type:    "success",
		})

		c.JSON(200, gin.H{
			"data":   symbol,
			"status": 200,
		})
	}
}
