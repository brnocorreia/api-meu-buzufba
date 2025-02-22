package controller

import (
	"net/http"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/controller/request"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/service"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"github.com/gin-gonic/gin"
)

func NewUserControllerInterface(
	serviceInterface service.UserService,
) UserControllerInterface {
	return &userControllerInterface{
		service: serviceInterface,
	}
}

type UserControllerInterface interface {
	FindUserByID(c *gin.Context)
	FindUserByEmail(c *gin.Context)

	UpdateUser(c *gin.Context)
}

type userControllerInterface struct {
	service service.UserService
}

func (uc *userControllerInterface) UpdateUser(c *gin.Context) {
	userDomainFromContext, ok := c.Get("user_domain")
	if !ok || userDomainFromContext == nil {
		restErr := rest_err.NewBadRequestError("user not found")
		c.JSON(restErr.Code, restErr)
		return
	}

	userId := userDomainFromContext.(domain.UserDomainInterface).GetID()

	var updateRequest request.UserUpdateRequest

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid request body")
		c.JSON(restErr.Code, restErr)
		return
	}

	err := uc.service.UpdateUser(userId, updateRequest)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (uc *userControllerInterface) FindUserByID(c *gin.Context) {
	restErr := rest_err.NewNotImplementedError("Endpoint not implemented yet")
	c.JSON(restErr.Code, restErr)
}

func (uc *userControllerInterface) FindUserByEmail(c *gin.Context) {
	restErr := rest_err.NewNotImplementedError("Endpoint not implemented yet")
	c.JSON(restErr.Code, restErr)
}
