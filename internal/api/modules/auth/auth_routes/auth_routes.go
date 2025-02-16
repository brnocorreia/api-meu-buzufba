package auth_routes

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/controller"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/domain/service"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitAuthRoutes(router *gin.RouterGroup, db *mongo.Database) {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authController := controller.NewAuthCoontrollerInterface(authService)

	authGroup := router.Group("auth")
	{
		authGroup.POST("/signin", authController.SignIn)
		authGroup.POST("/signup", authController.SignUp)
		authGroup.POST("/signout", authController.SignOut)
		authGroup.POST("/verify-email", authController.VerifyEmail)
	}
}
