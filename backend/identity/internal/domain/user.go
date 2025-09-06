package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserID string
type UserStatus int

const (
	StatusUnknown UserStatus = iota
	StatusPendingVerification
	StatusActive
	StatusSuspended
)

type User struct {
	ID              UserID
	Email           string
	PasswordHash    string
	Status          UserStatus
	TwoFactorSecret *string
	KYCLevel        int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewUser(email, plainPassword string) (*User, error) {
	if email == "" || plainPassword == "" {
		return nil, ErrEmailOrPasswordCantEmpty
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Email:        email,
		PasswordHash: string(hashPassword),
		Status:       StatusPendingVerification,
		KYCLevel:     0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) CheckPassword(plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return err
	}
	return nil
}

func (u *User) Enable2FA(secret string) {
	u.TwoFactorSecret = &secret
	u.UpdatedAt = time.Now()
}

func (u *User) ChangeStatus(status UserStatus) error {
	if !status.IsValid() {
		return ErrInvalidUserStatus
	}
	u.Status = status
	u.UpdatedAt = time.Now()

	return nil
}

func (s UserStatus) IsValid() bool {
	return s >= StatusPendingVerification && s <= StatusSuspended
}

func (s UserStatus) String() string {
	switch s {
	case StatusPendingVerification:
		return "PENDING_VERIFICATION"
	case StatusActive:
		return "ACTIVE"
	case StatusSuspended:
		return "SUSPENDED"
	default:
		return "UNKNOWN"
	}
}

// UserRegistered TODO Add event (events.go)
type UserRegistered struct {
	UserID string
	Email  string
}
