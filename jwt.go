package agin

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func CreateJwtToken(uid, secret string, liveHour time.Duration) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  uid,
		"exp":  time.Now().Add(liveHour).Unix(),
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseJwtToken(token string, secret string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	return claim.Claims.(jwt.MapClaims)["uid"].(string), nil
}