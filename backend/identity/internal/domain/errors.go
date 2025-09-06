package domain

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrEmailAlreadyExists       = errors.New("email already exists")
	ErrInvalidPassword          = errors.New("invalid password")
	ErrEmailOrPasswordCantEmpty = errors.New("email and password cannot be empty")
	ErrInvalidUserStatus        = errors.New("invalid user status")
	ErrInvalidToken             = errors.New("invalid token")
	ErrInternal                 = errors.New("internal error")
)
