package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/models"
)

const CookieJWTTImeSeconds = 24 * 60 * 60

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GetSignedToken(user *models.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       user.UserId.String(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.UserId.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	secretKey := config.GetSecretKey()

	signedString, err := token.SignedString([]byte(secretKey))

	return signedString, err

}

func ParseAndVerify(token string) (*MyJWTClaims, error) {
	secretKey := config.GetSecretKey()
	claims := &MyJWTClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) { return []byte(secretKey), nil })

	if err != nil {
		return nil, err
	}

	return claims, nil
}
