package model

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAll        = Role("all")
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Staff struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Role      Role      `json:"role" db:"role"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}
type RegisterStaffRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

type RegisterStaffResponse struct {
	AccessToken string `json:"token"`
}

type LoginStaffRequest struct {
	Username string `json:"username" `
	Password string `json:"password" validate:"required"`
}

type StaffWithToken struct {
	AccessToken string `json:"token"`
}
