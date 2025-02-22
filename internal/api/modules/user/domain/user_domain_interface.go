package domain

import (
	"time"
)

type UserDomainInterface interface {
	GetID() string
	GetFirstName() string
	GetLastName() string
	GetEmail() string
	GetPassword() string
	GetIsVerified() bool
	GetVerificationToken() string
	GetVerificationExpires() time.Time
	GetEmailVerifiedAt() *time.Time
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time

	SetID(id string)
	SetIsVerified(isVerified bool)
	SetVerificationToken(verificationToken string)
	SetVerificationExpires(verificationExpires time.Time)
	SetEmailVerifiedAt(emailVerifiedAt *time.Time)
	SetCreatedAt(createdAt time.Time)
	SetUpdatedAt(updatedAt time.Time)

	EncryptPassword()
	ComparePassword(password string) bool
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

func NewUserTokenDomain(
	id string,
	email string,
	firstName string,
	lastName string,
) UserDomainInterface {
	return &userDomain{
		id:        id,
		email:     email,
		firstName: firstName,
		lastName:  lastName,
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
