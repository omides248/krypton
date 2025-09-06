package auth

import (
	"errors"
	"fmt"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidAccessToken = errors.New("invalid access token")

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

type TokenManager struct {
	signingKey []byte
}

func NewTokenManager(signingKey string) *TokenManager {
	return &TokenManager{signingKey: []byte(signingKey)}
}

func (tm *TokenManager) newClaims(userID string, subject string, ttl time.Duration) Claims {
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}
}

func (tm *TokenManager) GenerateAccessToken(userID string) (string, error) {
	claims := tm.newClaims(userID, "access_token", time.Minute*15)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.signingKey)
}

func (tm *TokenManager) GenerateRefreshToken(userID string) (string, error) {
	claims := tm.newClaims(userID, "refresh_token", time.Hour*24*7)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.signingKey)
}

func (tm *TokenManager) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.signingKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (tm *TokenManager) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_ = c.Error(ErrInvalidAccessToken)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			_ = c.Error(ErrInvalidAccessToken)
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := tm.Validate(tokenString)
		if err != nil || claims.Subject != "access_token" {
			_ = c.Error(ErrInvalidAccessToken)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
