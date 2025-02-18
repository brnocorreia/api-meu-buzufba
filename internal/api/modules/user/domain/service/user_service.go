package service

import (
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

	UpdateUser(string, domain.UserDomainInterface) *rest_err.RestErr
}

func (us *userService) UpdateUser(
	userId string,
	userDomain domain.UserDomainInterface,
) *rest_err.RestErr {
	return nil
}

func (us *userService) FindUserByEmail(
	email string,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	return us.userRepository.FindUserByEmail(email)
}

func (us *userService) FindUserByID(
	id string,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	return nil, nil
}
