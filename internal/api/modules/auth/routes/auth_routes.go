package routes

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/controller"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/domain/service"

	authRepository "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/domain/repository"
	userRepository "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitAuthRoutes(router *gin.RouterGroup, db *mongo.Database) {
	userRepo := userRepository.NewUserRepository(db)
	authRepo := authRepository.NewAuthRepository(db)

	authService := service.NewAuthService(userRepo, authRepo)

	authController := controller.NewAuthCoontrollerInterface(authService)

	authGroup := router.Group("auth")
	{
		authGroup.POST("/signin", authController.SignIn)
		authGroup.POST("/signup", authController.SignUp)
		authGroup.POST("/signout", authController.SignOut)
		authGroup.POST("/email/verify", authController.VerifyEmail)
		authGroup.POST("/email/request-verification", service.VerifyTokenMiddleware, authController.RequestVerificationEmail)
	}
}
