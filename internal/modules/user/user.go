package user

import (
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/model"
	"github.com/brnocorreia/api-meu-buzufba/pkg/crypto"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/uid"
)

type user struct {
	id           string
	name         string
	username     string
	email        string
	password     string
	is_ufba      bool
	activated    bool
	activated_at *time.Time
	created_at   time.Time
	updated_at   time.Time
}

func NewFromModel(m model.User) *user {
	return &user{
		id:           m.ID,
		name:         m.Name,
		username:     m.Username,
		email:        m.Email,
		password:     m.Password,
		is_ufba:      m.IsUfba,
		activated:    m.Activated,
		activated_at: m.ActivatedAt,
		created_at:   m.CreatedAt,
		updated_at:   m.UpdatedAt,
	}
}

func New(name, username, email, pass string, isUfba bool) (*user, error) {
	hashedPass, err := crypto.HashPassword(pass)
	if err != nil {
		return nil, fault.New("failed to hash password", fault.WithError(err))
	}

	u := user{
		id:           uid.New("user"),
		name:         name,
		username:     username,
		email:        email,
		password:     hashedPass,
		is_ufba:      isUfba,
		activated:    false,
		activated_at: nil,
		created_at:   time.Now(),
		updated_at:   time.Now(),
	}

	if err := u.validate(); err != nil {
		return nil, fault.New(
			"failed to create user entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &u, nil
}

func (u *user) Activate() {
	u.activated = true
	now := time.Now()
	u.activated_at = &now
}

func (u *user) Model() model.User {
	return model.User{
		ID:          u.id,
		Name:        u.name,
		Username:    u.username,
		Email:       u.email,
		Password:    u.password,
		IsUfba:      u.is_ufba,
		Activated:   u.activated,
		ActivatedAt: u.activated_at,
		CreatedAt:   u.created_at,
		UpdatedAt:   u.updated_at,
	}
}

func (u *user) validate() error {
	if u.name == "" {
		return fault.New("user name is required")
	}
	if u.password == "" {
		return fault.New("password is required")
	}
	if u.email == "" {
		return fault.New("email is required")
	}
	if u.username == "" {
		return fault.New("username is required")
	}

	return nil
}
