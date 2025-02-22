package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthControllerInterface interface {
	Status(c *gin.Context)
}

type healthController struct {
	db    *mongo.Database
	redis *redis.Client
}

func NewHealthController(db *mongo.Database, redis *redis.Client) HealthControllerInterface {
	return &healthController{
		db:    db,
		redis: redis,
	}
}

func (hc *healthController) Status(c *gin.Context) {
	response := gin.H{
		"status": "healthy",
		"mongo":  "up",
		"redis":  "up",
	}

	if err := hc.db.Client().Ping(c, nil); err != nil {
		response["mongo"] = "down"
		response["status"] = "unhealthy"
	}

	if err := hc.redis.Ping(c).Err(); err != nil {
		response["redis"] = "down"
		response["status"] = "unhealthy"
	}

	if response["status"] == "healthy" {
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusServiceUnavailable, response)
}
