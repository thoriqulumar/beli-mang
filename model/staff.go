package model

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Staff struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username" db:"username"`
	Role      Role      `json:"role"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}
type RegisterStaffRequest struct {
	Username string `json:"username" `
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type RegisterStaffResponse struct {
	AccessToken string `json:"accessToken"`
}

type LoginStaffRequest struct {
	Username string `json:"username" `
	Password string `json:"password" validate:"required"`
}

type StaffWithToken struct {
	AccessToken string `json:"accessToken"`
}
