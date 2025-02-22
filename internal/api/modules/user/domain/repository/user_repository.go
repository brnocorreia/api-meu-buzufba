package repository

import (
	"context"
	"fmt"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/repository/entity"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/logger"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
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

	FindUserByID(
		id string,
	) (domain.UserDomainInterface, *rest_err.RestErr)

	FindUserByVerificationToken(
		verificationToken string,
	) (domain.UserDomainInterface, *rest_err.RestErr)
}

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
	collection := ur.db.Collection(USER_COLLECTION)

	value := entity.ConvertDomainToEntity(userDomain)

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return rest_err.NewBadRequestError("invalid user id")
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: value}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedUser entity.User
	err = collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return rest_err.NewNotFoundError("user not found")
		}
		logger.Error("error trying to update user", err, zap.String("user_id", userId))
		return rest_err.NewInternalServerError("Error trying to update user")
	}

	return nil
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

func (ur *userRepository) FindUserByID(
	id string,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	collection := ur.db.Collection(USER_COLLECTION)

	userEntity := &entity.User{}

	filter := bson.D{{Key: "_id", Value: id}}
	err := collection.FindOne(context.Background(), filter).Decode(userEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewNotFoundError("user not found with this id")
		}

		return nil, rest_err.NewInternalServerError("Error trying to find user by id")
	}

	return entity.ConvertEntityToDomain(userEntity), nil
}

func (ur *userRepository) FindUserByVerificationToken(
	verificationToken string,
) (domain.UserDomainInterface, *rest_err.RestErr) {
	collection := ur.db.Collection(USER_COLLECTION)

	userEntity := &entity.User{}

	filter := bson.D{{Key: "verification_token", Value: verificationToken}}
	err := collection.FindOne(context.Background(), filter).Decode(userEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewNotFoundError("user not found with this verification token")
		}

		return nil, rest_err.NewInternalServerError("Error trying to find user by verification token")
	}

	return entity.ConvertEntityToDomain(userEntity), nil
}
