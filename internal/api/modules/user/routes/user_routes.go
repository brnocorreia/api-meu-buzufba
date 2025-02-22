package routes

import (
	authService "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/domain/service"
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
		userGroup.PUT("/update", authService.VerifyTokenMiddleware, userController.UpdateUser)
	}
}
