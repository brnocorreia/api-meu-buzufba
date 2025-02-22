package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	FirstName           string             `bson:"first_name,omitempty"`
	LastName            string             `bson:"last_name,omitempty"`
	Email               string             `bson:"email,omitempty"`
	Password            string             `bson:"password,omitempty"`
	IsVerified          bool               `bson:"is_verified,omitempty" default:"false"`
	VerificationToken   string             `bson:"verification_token,omitempty"`
	VerificationExpires time.Time          `bson:"verification_expires,omitempty"`
	EmailVerifiedAt     *time.Time         `bson:"email_verified_at,omitempty"`
	CreatedAt           time.Time          `bson:"created_at,omitempty"`
	UpdatedAt           time.Time          `bson:"updated_at,omitempty"`
}
