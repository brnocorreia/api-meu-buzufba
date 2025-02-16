package controller

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"github.com/gin-gonic/gin"
)

func (uc *userControllerInterface) UpdateUser(c *gin.Context) {
	restErr := rest_err.NewNotImplementedError("Endpoint not implemented yet")
	c.JSON(restErr.Code, restErr)
}

func (uc *userControllerInterface) FindUserByID(c *gin.Context) {
	restErr := rest_err.NewNotImplementedError("Endpoint not implemented yet")
	c.JSON(restErr.Code, restErr)
}

func (uc *userControllerInterface) FindUserByEmail(c *gin.Context) {
	restErr := rest_err.NewNotImplementedError("Endpoint not implemented yet")
	c.JSON(restErr.Code, restErr)
}
