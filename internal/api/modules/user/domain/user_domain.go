package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userDomain struct {
	id         string
	firstName  string
	lastName   string
	email      string
	password   string
	isVerified bool
	createdAt  time.Time
	updatedAt  time.Time
}

func (ud *userDomain) GetID() string {
	return ud.id
}

func (ud *userDomain) GetFirstName() string {
	return ud.firstName
}

func (ud *userDomain) GetLastName() string {
	return ud.lastName
}

func (ud *userDomain) GetEmail() string {
	return ud.email
}

func (ud *userDomain) GetPassword() string {
	return ud.password
}

func (ud *userDomain) GetIsVerified() bool {
	return ud.isVerified
}

func (ud *userDomain) GetCreatedAt() time.Time {
	return ud.createdAt
}

func (ud *userDomain) GetUpdatedAt() time.Time {
	return ud.updatedAt
}

func (ud *userDomain) SetID(id string) {
	ud.id = id
}

func (ud *userDomain) SetIsVerified(isVerified bool) {
	ud.isVerified = isVerified
}

func (ud *userDomain) SetCreatedAt(createdAt time.Time) {
	ud.createdAt = createdAt
}

func (ud *userDomain) SetUpdatedAt(updatedAt time.Time) {
	ud.updatedAt = updatedAt
}

func (ud *userDomain) EncryptPassword() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ud.password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	ud.password = string(hashedPassword)
}

func (ud *userDomain) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(ud.password), []byte(password))
	return err == nil
}
