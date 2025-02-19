package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerifyToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	Token     string             `bson:"token,omitempty"`
	IsUsed    bool               `bson:"is_used,omitempty" default:"false"`
	ExpiresAt time.Time          `bson:"expires_at,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
}

func NewVerifyToken(
	userID primitive.ObjectID,
	expiresAt time.Time,
) *VerifyToken {
	return &VerifyToken{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Token:     uuid.New().String(),
		IsUsed:    false,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}
