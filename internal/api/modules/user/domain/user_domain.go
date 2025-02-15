package domain

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type userDomain struct {
	id        string
	firstName string
	lastName  string
	email     string
	password  string
	createdAt time.Time
	updatedAt time.Time
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

func (ud *userDomain) GetCreatedAt() time.Time {
	return ud.createdAt
}

func (ud *userDomain) GetUpdatedAt() time.Time {
	return ud.updatedAt
}

func (ud *userDomain) SetID(id string) {
	ud.id = id
}

func (ud *userDomain) SetCreatedAt(createdAt time.Time) {
	ud.createdAt = createdAt
}

func (ud *userDomain) SetUpdatedAt(updatedAt time.Time) {
	ud.updatedAt = updatedAt
}

func (ud *userDomain) EncryptPassword() {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(ud.password))
	ud.password = hex.EncodeToString(hash.Sum(nil))
}
