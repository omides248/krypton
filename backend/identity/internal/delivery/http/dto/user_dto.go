package dto

import "krypton/identity/internal/domain"

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:    string(user.ID),
		Email: user.Email,
	}
}
