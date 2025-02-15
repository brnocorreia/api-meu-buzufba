package controller

import (
	"net/http"

	user_request "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/controller/request"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/controller/response"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"github.com/gin-gonic/gin"
)

func (uc *userControllerInterface) CreateUser(c *gin.Context) {

	var userRequest user_request.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid request body")
		c.JSON(restErr.Code, restErr)
		return
	}

	domain := domain.NewUserDomain(
		userRequest.FirstName,
		userRequest.LastName,
		userRequest.Email,
		userRequest.Password,
	)

	user, err := uc.service.CreateUser(domain)
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

func (uc *userControllerInterface) UpdateUser(c *gin.Context) {

}

func (uc *userControllerInterface) Login(c *gin.Context) {

}

func (uc *userControllerInterface) FindUserByID(c *gin.Context) {

}

func (uc *userControllerInterface) FindUserByEmail(c *gin.Context) {

}
