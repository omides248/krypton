package error_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	StatusCode int         `json:"status_code"`
	ErrorType  string      `json:"error_type"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

type DomainErrorMapping struct {
	StatusCode int
	Message    string
}

const (
	SystemError     = "SYSTEM_ERROR"
	ValidationError = "VALIDATION_ERROR"
	DomainError     = "DOMAIN_ERROR"
)

func New(domainErrorMappings map[error]DomainErrorMapping, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0].Err

		var (
			statusCode = http.StatusInternalServerError
			errorType  = SystemError
			message    = "An unexpected error occurred."
			details    interface{}
		)

		var ozzoErrs ozzo.Errors
		var jsonErr *json.UnmarshalTypeError

		if errors.Is(err, io.EOF) {
			statusCode = http.StatusBadRequest
			errorType = ValidationError
			message = "Request body cannot be empty."

		} else if errors.As(err, &jsonErr) {
			statusCode = http.StatusBadRequest
			errorType = ValidationError
			message = fmt.Sprintf("Invalid type for field '%s'. Expected '%s'.", jsonErr.Field, jsonErr.Type.String())

		} else if errors.As(err, &ozzoErrs) {
			statusCode = http.StatusBadRequest
			errorType = ValidationError
			message = "Validation failed."

			mapped := make(map[string]string)
			for field, e := range ozzoErrs {
				jsonKey := toSnakeCase(field)
				mapped[jsonKey] = e.Error()
			}
			details = mapped
		} else {
			for domainErr, mapping := range domainErrorMappings {
				if errors.Is(err, domainErr) {
					statusCode = mapping.StatusCode
					errorType = DomainError
					message = mapping.Message
					break
				}
			}
		}

		logger.Error("handling http error",
			zap.Error(err),
			zap.Int("status_code", statusCode),
			zap.String("error_type", errorType),
		)

		c.AbortWithStatusJSON(statusCode, ErrorResponse{
			StatusCode: statusCode,
			ErrorType:  errorType,
			Message:    message,
			Details:    details,
		})
	}
}

func toSnakeCase(s string) string {
	var sb strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
