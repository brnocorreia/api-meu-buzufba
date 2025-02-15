package repository

import (
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(
	db *mongo.Database,
) UserRepository {
	return &userRepository{
		db,
	}
}

type userRepository struct {
	db *mongo.Database
}

type UserRepository interface {
	CreateUser(
		userDomain domain.UserDomainInterface,
	) (domain.UserDomainInterface, *rest_err.RestErr)

	UpdateUser(
		userId string,
		userDomain domain.UserDomainInterface,
	) *rest_err.RestErr

	FindUserByEmail(
		email string,
	) (domain.UserDomainInterface, *rest_err.RestErr)

	FindUserByEmailAndPassword(
		email string,
		password string,
	) (domain.UserDomainInterface, *rest_err.RestErr)

	FindUserByID(
		id string,
	) (domain.UserDomainInterface, *rest_err.RestErr)
}
