package repository

import (
	"context"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/auth/domain/repository/entity"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAuthRepository(
	db *mongo.Database,
) AuthRepository {
	return &authRepository{
		db,
	}
}

type authRepository struct {
	db *mongo.Database
}

type AuthRepository interface {
	CreateVerifyEmailToken(
		userID primitive.ObjectID,
		expiresAt time.Time,
	) (*entity.VerifyToken, *rest_err.RestErr)

	FindVerifyEmailTokenByUserID(
		userID primitive.ObjectID,
	) (*entity.VerifyToken, *rest_err.RestErr)

	FindVerifyEmailTokenByID(
		id primitive.ObjectID,
	) (*entity.VerifyToken, *rest_err.RestErr)

	UpdateVerifyEmailToken(
		id primitive.ObjectID,
		verifyToken entity.VerifyToken,
	) (*entity.VerifyToken, *rest_err.RestErr)
}

const (
	VERIFY_EMAIL_TOKEN_COLLECTION = "verify_email_tokens"
)

func (ar *authRepository) CreateVerifyEmailToken(
	userID primitive.ObjectID,
	expiresAt time.Time,
) (*entity.VerifyToken, *rest_err.RestErr) {
	verifyToken := entity.NewVerifyToken(userID, expiresAt)

	_, err := ar.db.Collection(VERIFY_EMAIL_TOKEN_COLLECTION).InsertOne(context.Background(), verifyToken)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error creating verify email token")
	}

	return verifyToken, nil
}

func (ar *authRepository) FindVerifyEmailTokenByUserID(
	userID primitive.ObjectID,
) (*entity.VerifyToken, *rest_err.RestErr) {
	filter := bson.M{"user_id": userID}

	var verifyToken entity.VerifyToken
	err := ar.db.Collection(VERIFY_EMAIL_TOKEN_COLLECTION).FindOne(context.Background(), filter).Decode(&verifyToken)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error finding verify email token")
	}

	return &verifyToken, nil
}

func (ar *authRepository) FindVerifyEmailTokenByID(
	id primitive.ObjectID,
) (*entity.VerifyToken, *rest_err.RestErr) {
	filter := bson.M{"_id": id}

	var verifyToken entity.VerifyToken
	err := ar.db.Collection(VERIFY_EMAIL_TOKEN_COLLECTION).FindOne(context.Background(), filter).Decode(&verifyToken)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error finding verify email token")
	}

	return &verifyToken, nil
}

func (ar *authRepository) UpdateVerifyEmailToken(
	id primitive.ObjectID,
	verifyToken entity.VerifyToken,
) (*entity.VerifyToken, *rest_err.RestErr) {
	filter := bson.M{"_id": id}

	_, err := ar.db.Collection(VERIFY_EMAIL_TOKEN_COLLECTION).UpdateOne(context.Background(), filter, bson.M{"$set": verifyToken})
	if err != nil {
		return nil, rest_err.NewInternalServerError("error updating verify email token")
	}

	return &verifyToken, nil
}
