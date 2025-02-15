package service

import (
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
)

func (us *userService) CreateUser(
	userDomain domain.UserDomainInterface,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	user, _ := us.FindUserByEmail(userDomain.GetEmail())
	if user != nil {
		return nil, rest_err.NewBadRequestError("Email already registered in another account")
	}

	userDomain.EncryptPassword()
	userDomain.SetCreatedAt(time.Now())
	userDomain.SetUpdatedAt(time.Now())
	user, err := us.userRepository.CreateUser(userDomain)
	if err != nil {
		return nil, err
	}

	return user, nil
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

func (us *userService) Login(
	userDomain domain.UserDomainInterface,
) (domain.UserDomainInterface, string, *rest_err.RestErr) {
	return nil, "", nil
}
