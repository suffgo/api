package infrastructure

import (
	"suffgo/cmd/config"
	valueobjects "suffgo/internal/user/domain/valueObjects"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var secretKey = []byte(config.SecretKey)

func createToken(username valueobjects.UserName, ip, agent string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	sessionID := uuid.New().String()

	claims["username"] = username.Username
	claims["session_id"] = sessionID
	claims["ip"] = ip
	claims["user_agent"] = agent
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return t, nil
}
