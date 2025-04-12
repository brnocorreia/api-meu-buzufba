package model

import "time"

type Session struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	IP           string    `db:"ip_address"`
	Agent        string    `db:"agent"`
	RefreshToken string    `db:"refresh_token"`
	Active       bool      `db:"active"`
	Expires      time.Time `db:"expires"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type User struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Username    string    `db:"username"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Activated   bool      `db:"activated"`
	ActivatedAt time.Time `db:"activated_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
