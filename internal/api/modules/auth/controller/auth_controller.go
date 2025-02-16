package controller

import (
	"net/http"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/controller/request"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/domain/service"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/controller/response"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"github.com/gin-gonic/gin"
)

func NewAuthCoontrollerInterface(service service.AuthService) AuthControllerInterface {
	return &authControllerInterface{
		service: service,
	}
}

type AuthControllerInterface interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
	SignOut(c *gin.Context)
	VerifyEmail(c *gin.Context)
}

type authControllerInterface struct {
	service service.AuthService
}

func (ac *authControllerInterface) SignIn(c *gin.Context) {
	restErr := rest_err.NewNotImplementedError("Endpoint not implemented yet")
	c.JSON(restErr.Code, restErr)
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

	user, err := ac.service.SignUp(domain)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	response := response.UserResponse{
		ID:        user.GetID(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Email:     user.GetEmail(),
		CreatedAt: user.GetCreatedAt(),
		UpdatedAt: user.GetUpdatedAt(),
	}

	c.JSON(http.StatusCreated, response)
}

func (ac *authControllerInterface) SignOut(c *gin.Context) {
	restErr := rest_err.NewNotImplementedError("Endpoint not implemented yet")
	c.JSON(restErr.Code, restErr)
}

func (ac *authControllerInterface) VerifyEmail(c *gin.Context) {
	restErr := rest_err.NewNotImplementedError("Endpoint not implemented yet")
	c.JSON(restErr.Code, restErr)
}
