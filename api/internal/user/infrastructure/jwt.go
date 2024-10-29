package infrastructure

import (
	"suffgo/cmd/config"
	valueobjects "suffgo/internal/user/domain/valueObjects"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte(config.SecretKey) 

func createToken(username valueobjects.UserName) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return t, nil
}
