package service

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/controller/request"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/repository"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
)

func NewUserService(
	userRepository repository.UserRepository,
) UserService {
	return &userService{userRepository}
}

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	FindUserByEmail(
		email string,
	) (domain.UserDomainInterface, *rest_err.RestErr)

	FindUserByID(
		id string,
	) (domain.UserDomainInterface, *rest_err.RestErr)

	UpdateUser(string, request.UserUpdateRequest) *rest_err.RestErr
}

func (us *userService) UpdateUser(
	userId string,
	updateUserRequest request.UserUpdateRequest,
) *rest_err.RestErr {
	userDomain, err := us.FindUserByID(userId)
	if err != nil {
		return err
	}

	userDomain.SetFirstName(updateUserRequest.FirstName)
	userDomain.SetLastName(updateUserRequest.LastName)

	return us.userRepository.UpdateUser(userId, userDomain)
}

func (us *userService) FindUserByEmail(
	email string,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	return us.userRepository.FindUserByEmail(email)
}

func (us *userService) FindUserByID(
	id string,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	return us.userRepository.FindUserByID(id)
}
