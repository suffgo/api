package infrastructure

import (
	"suffgo/cmd/config"
	valueobjects "suffgo/internal/user/domain/valueObjects"
	"time"
	sv "suffgo/internal/shared/domain/valueObjects"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var secretKey = []byte(config.SecretKey)

func createToken(username valueobjects.UserName, ip, agent string, userID sv.ID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	sessionID := uuid.New().String()

	claims["user_id"] = userID.Id
	claims["username"] = username.Username
	claims["session_id"] = sessionID
	claims["ip"] = ip
	claims["user_agent"] = agent
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return t, nil
}
