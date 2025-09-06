package dto

import ozzo "github.com/go-ozzo/ozzo-validation/v4"
import "github.com/go-ozzo/ozzo-validation/v4/is"

// --------------------- Register ---------------------

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Email, ozzo.Required, is.Email),
		ozzo.Field(&r.Password, ozzo.Required, ozzo.Length(6, 100)),
	)
}

// --------------------- Login ---------------------

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Email, ozzo.Required, is.Email),
		ozzo.Field(&r.Password, ozzo.Required),
	)
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token Token        `json:"token"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r Token) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.RefreshToken, ozzo.Required),
	)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.RefreshToken, ozzo.Required),
	)
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
