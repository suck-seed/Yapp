package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/models"
)

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GetSignedToken(user *models.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       user.ID.String(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	secretKey := config.GetSecretKey()

	signedString, err := token.SignedString([]byte(secretKey))

	return signedString, err

}
