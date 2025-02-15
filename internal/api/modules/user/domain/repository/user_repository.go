package repository

import (
	"context"
	"fmt"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/repository/entity"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	USER_COLLECTION = "users"
)

func (ur *userRepository) CreateUser(
	userDomain domain.UserDomainInterface,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	value := entity.ConvertDomainToEntity(userDomain)

	collection := ur.db.Collection(USER_COLLECTION)

	result, err := collection.InsertOne(context.Background(), value)
	if err != nil {
		return nil, rest_err.NewInternalServerError(err.Error())
	}

	value.ID = result.InsertedID.(primitive.ObjectID)

	return entity.ConvertEntityToDomain(value), nil
}

func (ur *userRepository) UpdateUser(
	userId string,
	userDomain domain.UserDomainInterface,
) *rest_err.RestErr {
	return rest_err.NewNotImplementedError("method not implemented")
}

func (ur *userRepository) FindUserByEmail(
	email string,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	collection := ur.db.Collection(USER_COLLECTION)

	userEntity := &entity.User{}

	filter := bson.D{{Key: "email", Value: email}}
	err := collection.FindOne(context.Background(), filter).Decode(userEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			errorMessage := fmt.Sprintf("user not found with this email: %s", email)
			return nil, rest_err.NewNotFoundError(errorMessage)
		}

		return nil, rest_err.NewInternalServerError("Error trying to find user by email")
	}

	return entity.ConvertEntityToDomain(userEntity), nil
}

func (ur *userRepository) FindUserByEmailAndPassword(
	email string,
	password string,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	return nil, rest_err.NewNotImplementedError("method not implemented")
}

func (ur *userRepository) FindUserByID(
	id string,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	return nil, rest_err.NewNotImplementedError("method not implemented")
}
