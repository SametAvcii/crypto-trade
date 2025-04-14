package routes

import (
	"net/http"

	ctlog "github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/domains/exchange"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/gin-gonic/gin"
)

func ExchangeRoutes(r *gin.RouterGroup, s exchange.Service) {
	r.POST("/", AddExchange(s))
	r.PUT("/:id", UpdateExchange(s))
	r.GET("/:id", GetExchangeById(s))
	r.DELETE("/:id", DeleteExchange(s))
	r.GET("/", GetAllExchanges(s))
}

// AddExchange godoc
// @Summary Add exchange
// @Description Add a new exchange
// @Tags exchanges
// @Accept json
// @Produce json
// @Param exchange body dtos.AddExchangeReq true "Exchange information"
// @Success 201 {object} dtos.AddExchangeRes "Successfully created exchange"
// @Failure 400 {object} map[string]any
// @Router /exchanges [post]
func AddExchange(s exchange.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dtos.AddExchangeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error":  err.Error(),
				"status": 400,
			})
			return
		}

		res, err := s.AddExchange(c, req)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Add Exchange Error",
				Message: "Add Exchange err: " + err.Error(),
				Entity:  "exchange",
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
			Entity:  "exchange",
			Type:    "success",
		})

		c.JSON(201, gin.H{
			"data":   res,
			"status": 201,
		})
	}
}

// UpdateExchange godoc
// @Summary Update an exchange
// @Description Updates an existing exchange with the provided information
// @Tags exchanges
// @Accept json
// @Produce json
// @Param id path string true "Exchange ID"
// @Param exchange body dtos.UpdateExchangeReq true "Exchange update information"
// @Success 200 {object} dtos.UpdateExchangeRes "Successfully updated exchange"
// @Failure 400 {object} map[string]any
// @Router /exchanges/{id} [put]
func UpdateExchange(s exchange.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dtos.UpdateExchangeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error(), "status": 400})
			return
		}
		req.ID = c.Param("id")
		res, err := s.Update(c, req)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Update Exchange Error",
				Message: "Update exchange err: " + err.Error(),
				Entity:  "exchange",
				Type:    "error",
			})
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}
		ctlog.CreateLog(&entities.Log{
			Title:   "Update Exchange",
			Message: "Update exchange success: " + res.Name,
			Entity:  "exchange",
			Type:    "success",
		})

		c.JSON(200, gin.H{"data": res, "status": 200})
	}
}

// GetExchangeById godoc
// @Summary Get exchange by ID
// @Description Retrieves an exchange by its ID
// @Tags exchanges
// @Accept json
// @Produce json
// @Param id path string true "Exchange ID"
// @Success 200 {object} dtos.GetExchangeRes
// @Failure 400 {object} map[string]any
// @Router /exchanges/{id} [get]
func GetExchangeById(s exchange.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		res, err := s.GetById(c, id)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Get Exchange Error",
				Message: "Get exchange err: " + err.Error(),
				Entity:  "exchange",
				Type:    "error",
			})
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}
		ctlog.CreateLog(&entities.Log{
			Title:   "Get Exchange",
			Message: "Get exchange success: " + res.Name,
			Entity:  "exchange",
			Type:    "success",
		})

		c.JSON(200, gin.H{"data": res, "status": 200})
	}
}

// DeleteExchange godoc
// @Summary Delete an exchange
// @Description Deletes an exchange by its ID
// @Tags exchanges
// @Accept json
// @Produce json
// @Param id path string true "Exchange ID"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /exchanges/{id} [delete]
func DeleteExchange(s exchange.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := s.Delete(c, id); err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Delete Exchange Error",
				Message: "Delete exchange err: " + err.Error(),
				Entity:  "exchange",
				Type:    "error",
			})
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}
		ctlog.CreateLog(&entities.Log{
			Title:   "Delete Exchange",
			Message: "Delete exchange success",
			Entity:  "exchange",
			Type:    "success",
		})

		c.JSON(200, gin.H{"message": "Exchange deleted successfully", "status": 200})
	}
}

// GetAllExchanges godoc
// @Summary Get all exchanges
// @Description Retrieves all exchanges
// @Tags exchanges
// @Accept json
// @Produce json
// @Success 200 {array} dtos.GetExchangeRes
// @Failure 400 {object} map[string]any
// @Router /exchanges [get]
func GetAllExchanges(s exchange.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		res, err := s.GetAll(c)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Get All Exchanges Error",
				Message: "Get all exchanges err: " + err.Error(),
				Entity:  "exchange",
				Type:    "error",
			})

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}

		ctlog.CreateLog(&entities.Log{
			Title:   "Get All Exchanges",
			Message: "Get all exchanges success",
			Entity:  "exchange",
			Type:    "success",
		})
		c.JSON(200, gin.H{"data": res, "status": 200})
	}
}
