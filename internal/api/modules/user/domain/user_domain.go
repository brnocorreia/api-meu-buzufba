package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userDomain struct {
	id                  string
	firstName           string
	lastName            string
	email               string
	password            string
	isVerified          bool
	verificationToken   string
	verificationExpires time.Time
	emailVerifiedAt     *time.Time
	createdAt           time.Time
	updatedAt           time.Time
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

func (ud *userDomain) GetVerificationToken() string {
	return ud.verificationToken
}

func (ud *userDomain) GetVerificationExpires() time.Time {
	return ud.verificationExpires
}

func (ud *userDomain) GetEmailVerifiedAt() *time.Time {
	return ud.emailVerifiedAt
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

func (ud *userDomain) SetFirstName(firstName string) {
	ud.firstName = firstName
}

func (ud *userDomain) SetLastName(lastName string) {
	ud.lastName = lastName
}

func (ud *userDomain) SetIsVerified(isVerified bool) {
	ud.isVerified = isVerified
}

func (ud *userDomain) SetVerificationToken(verificationToken string) {
	ud.verificationToken = verificationToken
}

func (ud *userDomain) SetVerificationExpires(verificationExpires time.Time) {
	ud.verificationExpires = verificationExpires
}

func (ud *userDomain) SetEmailVerifiedAt(emailVerifiedAt *time.Time) {
	ud.emailVerifiedAt = emailVerifiedAt
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

func (ud *userDomain) GenerateVerificationInfo() {
	ud.verificationToken = uuid.New().String()
	ud.verificationExpires = time.Now().Add(time.Hour * 24 * 3)
	ud.emailVerifiedAt = nil
	ud.isVerified = false
}
