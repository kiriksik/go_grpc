package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kiriksik/go_grpc/internal/domain/models"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	JWTToken, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}
	return JWTToken, nil
}
