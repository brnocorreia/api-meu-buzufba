package routes

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/auth_routes"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/user_routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoutes(r *gin.Engine, db *mongo.Database) {
	v1 := r.Group("/v1")
	user_routes.InitUserRoutes(v1, db)
	auth_routes.InitAuthRoutes(v1, db)
}
