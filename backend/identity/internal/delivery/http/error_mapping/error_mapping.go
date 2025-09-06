package error_mapping

import (
	"krypton/identity/internal/domain"
	"krypton/pkg/auth"
	"krypton/pkg/gin/error_handler"
	"net/http"
)

func GetDomainErrorMappings() map[error]error_handler.DomainErrorMapping {
	return map[error]error_handler.DomainErrorMapping{
		domain.ErrUserNotFound:       {StatusCode: http.StatusNotFound, Message: "user not found"},
		domain.ErrEmailAlreadyExists: {StatusCode: http.StatusConflict, Message: "email already exists"},
		domain.ErrInvalidPassword:    {StatusCode: http.StatusUnauthorized, Message: "invalid credentials"},
		domain.ErrInvalidToken:       {StatusCode: http.StatusUnauthorized, Message: domain.ErrInvalidToken.Error()},
		auth.ErrInvalidAccessToken:   {StatusCode: http.StatusUnauthorized, Message: "invalid access token"},
	}
}
