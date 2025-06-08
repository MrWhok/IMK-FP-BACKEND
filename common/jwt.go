package common

import (
	"strconv"
	"time"

	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(username string, roles []map[string]interface{}, config configuration.Config) (string, time.Time) {
	jwtSecret := config.Get("JWT_SECRET_KEY")
	jwtExpired, err := strconv.Atoi(config.Get("JWT_EXPIRE_MINUTES_COUNT"))
	exception.PanicLogging(err)

	expirationTime := time.Now().Add(time.Minute * time.Duration(jwtExpired))

	claims := jwt.MapClaims{
		"username": username,
		"roles":    roles,
		"exp":      expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(jwtSecret))
	exception.PanicLogging(err)

	return tokenSigned, expirationTime
}
