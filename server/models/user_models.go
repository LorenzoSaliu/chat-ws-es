package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id" validate:"required,min=2,max=255"`
	FirstName    string    `json:"firs_name" validate:"required,min=2,max=255"`
	LastName     string    `json:"last_name" validate:"required,min=2,max=255"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,min=6"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
