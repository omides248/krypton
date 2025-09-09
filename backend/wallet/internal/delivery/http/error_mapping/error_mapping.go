package error_mapping

import (
	"krypton/pkg/auth"
	"krypton/pkg/contextkeys"
	"krypton/pkg/gin/error_handler"
	"krypton/wallet/internal/domain"
	"net/http"
)

func GetDomainErrorMappings() map[error]error_handler.DomainErrorMapping {
	return map[error]error_handler.DomainErrorMapping{
		domain.ErrAssetNotFound:       {StatusCode: http.StatusNotFound, Message: "asset not found"},
		domain.ErrAccountNotFound:     {StatusCode: http.StatusNotFound, Message: "account not found"},
		auth.ErrInvalidAccessToken:    {StatusCode: http.StatusUnauthorized, Message: "invalid access token"},
		contextkeys.ErrUserIDNotFound: {StatusCode: http.StatusUnauthorized, Message: "user id not found in context"},
	}
}
