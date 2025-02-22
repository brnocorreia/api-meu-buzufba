package routes

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/health/controller"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitHealthRoutes(router *gin.RouterGroup, db *mongo.Database, redis *redis.Client) {
	healthController := controller.NewHealthController(db, redis)
	router.GET("/health", healthController.Status)
}
