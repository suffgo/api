package infrastructure

import (
	"suffgo/cmd/config"
	sv "suffgo/internal/shared/domain/valueObjects"
	valueobjects "suffgo/internal/user/domain/valueObjects"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte(config.SecretKey)

func createToken(username valueobjects.UserName, ip, agent string, id sv.ID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id.Id
	claims["username"] = username.Username
	claims["ip"] = ip
	claims["user_agent"] = agent
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return t, nil
}
