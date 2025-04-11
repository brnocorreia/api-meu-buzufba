package user

import (
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/model"
	"github.com/brnocorreia/api-meu-buzufba/pkg/crypto"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/uid"
)

type user struct {
	id         string
	name       string
	username   string
	email      string
	password   string
	avatar_url *string
	enabled    bool
	locked     bool
	created    time.Time
	updated    time.Time
}

func NewFromModel(m model.User) *user {
	return &user{
		id:         m.ID,
		name:       m.Name,
		username:   m.Username,
		email:      m.Email,
		password:   m.Password,
		avatar_url: m.AvatarURL,
		enabled:    m.Enabled,
		locked:     m.Locked,
		created:    m.CreatedAt,
		updated:    m.UpdatedAt,
	}
}

func New(name, username, email, pass string) (*user, error) {
	hashedPass, err := crypto.HashPassword(pass)
	if err != nil {
		return nil, fault.New("failed to hash password", fault.WithError(err))
	}

	u := user{
		id:         uid.New("user"),
		name:       name,
		username:   username,
		email:      email,
		password:   hashedPass,
		avatar_url: nil,
		enabled:    false,
		locked:     false,
		created:    time.Now(),
		updated:    time.Now(),
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

func (u *user) Enable() {
	u.enabled = true
	u.updated = time.Now()
}

func (u *user) ToModel() model.User {
	return model.User{
		ID:        u.id,
		Name:      u.name,
		Username:  u.username,
		Email:     u.email,
		Password:  u.password,
		AvatarURL: u.avatar_url,
		Enabled:   u.enabled,
		Locked:    u.locked,
		CreatedAt: u.created,
		UpdatedAt: u.updated,
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
