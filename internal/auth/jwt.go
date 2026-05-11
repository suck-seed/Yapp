package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/models"
)

const (
	AccessTokenTTL  = 60 * time.Minute
	RefreshTokenTTL = 30 * 24 * time.Hour
)

type JWTPayload struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GetSignedToken(user *models.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTPayload{
		ID:       user.ID.String(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "yapp",
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	secretKey := config.GetSecretKey()
	signedString, err := token.SignedString([]byte(secretKey))
	return signedString, err

}

func ParseAndVerify(tokenString string) (*JWTPayload, error) {
	secretKey := config.GetSecretKey()
	claims := &JWTPayload{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {

		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(secretKey), nil

	})
	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
