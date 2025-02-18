package routes

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/controller"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/repository"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitUserRoutes(router *gin.RouterGroup, db *mongo.Database) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserControllerInterface(userService)

	userGroup := router.Group("users")
	{
		userGroup.GET("/:userId", userController.FindUserByID)
		userGroup.GET("/email/:userEmail", userController.FindUserByEmail)
		userGroup.PUT("/:userId", userController.UpdateUser)
	}
}
