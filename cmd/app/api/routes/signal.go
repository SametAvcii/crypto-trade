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
	r.PUT("/interval/:id", UpdateSignalInterval(s))
	r.DELETE("/interval/:id", DeleteSignalInterval(s))
	r.GET("/interval/:id", GetSignalIntervalByID(s))
	r.GET("/interval", GetAllSignalIntervals(s))

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

		res, err := s.AddSignalIntervals(c, req)
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

// @Summary Update Signal Interval
// @Description Update Signal Interval Route
// @Tags Signal Endpoints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Signal ID"
// @Param payload body dtos.UpdateSignalIntervalReq true "Update Signal Request"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /signal/interval/{id} [PUT]
func UpdateSignalInterval(s signal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req dtos.UpdateSignalIntervalReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error(), "status": 400})
			return
		}
		req.ID = id

		res, err := s.UpdateSignalIntervals(c, req)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Update Signal Interval Error",
				Message: "Update Signal Interval err: " + err.Error(),
				Entity:  "signal",
				Type:    "error",
			})
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}

		ctlog.CreateLog(&entities.Log{
			Title:   "Update Signal Interval",
			Message: "Update Signal Interval success: " + res.Symbol,
			Entity:  "signal",
			Type:    "success",
		})
		c.JSON(200, gin.H{"data": res, "status": 200})
	}
}

// @Summary Delete Signal Interval
// @Description Delete Signal Interval Route
// @Tags Signal Endpoints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Signal ID"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /signal/interval/{id} [DELETE]
func DeleteSignalInterval(s signal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := s.DeleteSignalIntervals(c, id)
		if err != nil {
			ctlog.CreateLog(&entities.Log{
				Title:   "Delete Signal Interval Error",
				Message: "Delete Signal Interval err: " + err.Error(),
				Entity:  "signal",
				Type:    "error",
			})
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}

		ctlog.CreateLog(&entities.Log{
			Title:   "Delete Signal Interval",
			Message: "Delete Signal Interval success",
			Entity:  "signal",
			Type:    "success",
		})
		c.JSON(200, gin.H{"message": "Successfully deleted", "status": 200})
	}
}

// @Summary Get Signal Interval by ID
// @Description Get Signal Interval by ID Route
// @Tags Signal Endpoints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Signal ID"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /signal/interval/{id} [GET]
func GetSignalIntervalByID(s signal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		res, err := s.GetSignalIntervalById(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}
		c.JSON(200, gin.H{"data": res, "status": 200})
	}
}

// @Summary Get All Signal Intervals
// @Description Get All Signal Intervals Route
// @Tags Signal Endpoints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /signal/interval [GET]
func GetAllSignalIntervals(s signal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		res, err := s.GetAllSignalIntervals(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}
		c.JSON(200, gin.H{"data": res, "status": 200})
	}
}
