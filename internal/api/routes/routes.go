package routes

import (
	authRoutes "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/routes"
	userRoutes "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/routes"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoutes(r *gin.Engine, db *mongo.Database, redis *redis.Client) {
	v1 := r.Group("/v1")
	userRoutes.InitUserRoutes(v1, db)
	authRoutes.InitAuthRoutes(v1, db, redis)
}
