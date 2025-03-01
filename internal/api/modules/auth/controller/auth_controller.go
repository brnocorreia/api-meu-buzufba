package controller

import (
	"net/http"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/controller/request"
	authResponse "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/controller/response"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/domain/service"
	userResponse "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/controller/response"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rate_limiter"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"github.com/gin-gonic/gin"
)

func NewAuthCoontrollerInterface(service service.AuthService, rateLimiter *rate_limiter.RateLimiter) AuthControllerInterface {
	return &authControllerInterface{
		service:     service,
		rateLimiter: rateLimiter,
	}
}

type AuthControllerInterface interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
	SignOut(c *gin.Context)
	VerifyEmail(c *gin.Context)
	RequestVerificationEmail(c *gin.Context)
}

type authControllerInterface struct {
	service     service.AuthService
	rateLimiter *rate_limiter.RateLimiter
}

func (ac *authControllerInterface) SignIn(c *gin.Context) {
	var signInRequest request.SignInRequest

	if err := c.ShouldBindJSON(&signInRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid request body")
		c.JSON(restErr.Code, restErr)
		return
	}

	token, err := ac.service.SignIn(signInRequest.Email, signInRequest.Password)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Header("Authorization", token)

	response := authResponse.SignInResponse{
		AccessToken: token,
	}

	c.JSON(http.StatusOK, response)
}

func (ac *authControllerInterface) SignUp(c *gin.Context) {
	var signUpRequest request.SignUpRequest

	if err := c.ShouldBindJSON(&signUpRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid request body")
		c.JSON(restErr.Code, restErr)
		return
	}

	domain := domain.NewUserDomain(
		signUpRequest.FirstName,
		signUpRequest.LastName,
		signUpRequest.Email,
		signUpRequest.Password,
	)
	domain.SetIsVerified(false)

	user, err := ac.service.SignUp(domain)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	response := userResponse.UserResponse{
		ID:         user.GetID(),
		FirstName:  user.GetFirstName(),
		LastName:   user.GetLastName(),
		Email:      user.GetEmail(),
		IsVerified: user.GetIsVerified(),
		CreatedAt:  user.GetCreatedAt(),
		UpdatedAt:  user.GetUpdatedAt(),
	}

	c.JSON(http.StatusCreated, response)
}

func (ac *authControllerInterface) SignOut(c *gin.Context) {
	err := ac.service.SignOut(c)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (ac *authControllerInterface) VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		restErr := rest_err.NewBadRequestError("verification token is required")
		c.JSON(restErr.Code, restErr)
		return
	}

	err := ac.service.VerifyEmail(token)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (ac *authControllerInterface) RequestVerificationEmail(c *gin.Context) {
	userDomainFromContext, ok := c.Get("user_domain")
	if !ok || userDomainFromContext == nil {
		restErr := rest_err.NewBadRequestError("user not found")
		c.JSON(restErr.Code, restErr)
		return
	}

	userDomain := userDomainFromContext.(domain.UserDomainInterface)

	allowed, redisErr := ac.rateLimiter.IsAllowed(c.Request.Context(), userDomain.GetEmail())
	if redisErr != nil {
		restErr := rest_err.NewInternalServerError("error checking rate limit")
		c.JSON(restErr.Code, restErr)
		return
	}

	if !allowed {
		restErr := rest_err.NewTooManyRequestsError("too many requests")
		c.JSON(restErr.Code, restErr)
		return
	}

	err := ac.service.RequestVerificationEmail(userDomain.GetEmail())
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
