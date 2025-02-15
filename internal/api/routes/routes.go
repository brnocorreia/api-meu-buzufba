package routes

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoutes(r *gin.Engine, db *mongo.Database) {
	v1 := r.Group("api/v1")
	routes.InitUserRoutes(v1, db)

}
