package domain

import (
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
)

type UserDomainInterface interface {
	GetID() string
	GetFirstName() string
	GetLastName() string
	GetEmail() string
	GetPassword() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time

	SetID(id string)
	SetCreatedAt(createdAt time.Time)
	SetUpdatedAt(updatedAt time.Time)

	EncryptPassword()
	GenerateToken() (string, *rest_err.RestErr)
}

func NewUserDomain(
	firstName string,
	lastName string,
	email string,
	password string,

) UserDomainInterface {
	return &userDomain{
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		password:  password,
	}
}

func NewUserUpdateDomain(
	firstName string,
	lastName string,
	password string,
) UserDomainInterface {
	return &userDomain{
		firstName: firstName,
		lastName:  lastName,
		password:  password,
	}
}

func NewUserLoginDomain(
	email string,
	password string,
) UserDomainInterface {
	return &userDomain{
		email:    email,
		password: password,
	}
}
