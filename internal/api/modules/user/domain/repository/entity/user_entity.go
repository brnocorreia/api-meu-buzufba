package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	FirstName           string             `bson:"first_name"`
	LastName            string             `bson:"last_name"`
	Email               string             `bson:"email"`
	Password            string             `bson:"password"`
	IsVerified          bool               `bson:"is_verified"`
	VerificationToken   string             `bson:"verification_token"`
	VerificationExpires time.Time          `bson:"verification_expires"`
	EmailVerifiedAt     *time.Time         `bson:"email_verified_at,omitempty"`
	CreatedAt           time.Time          `bson:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at"`
}
