package routes

import (
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/controller"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/domain/service"
	userRepository "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/repository"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rate_limiter"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitAuthRoutes(router *gin.RouterGroup, db *mongo.Database, redis *redis.Client) {
	userRepo := userRepository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo)

	rateLimiter := rate_limiter.NewRateLimiter(redis, 3, time.Hour*2)

	authController := controller.NewAuthCoontrollerInterface(authService, rateLimiter)

	authGroup := router.Group("auth")
	{
		authGroup.POST("/signin", authController.SignIn)
		authGroup.POST("/signup", authController.SignUp)
		authGroup.POST("/signout", authController.SignOut)
		authGroup.POST("/email/verify", authController.VerifyEmail)
		authGroup.POST("/email/request-verification", service.VerifyTokenMiddleware, authController.RequestVerificationEmail)
	}
}
