package controller

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/service"
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

	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	Login(c *gin.Context)
}

type userControllerInterface struct {
	service service.UserService
}
