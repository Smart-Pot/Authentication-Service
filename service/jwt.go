package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = "jwtsecret"

type jwtClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

type jwtService struct{}

func (j jwtService) tokenize(userID string) (string, error) {

	expirationTime := time.Now().Add(24 * 30 * time.Hour)
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["userId"] = userID
	atClaims["exp"] = expirationTime.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j jwtService) verify(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token can not be resolved")
	}
	return claims, nil
}

func (j jwtService) getUserID(tokenStr string) (string, error) {
	claims, err := j.verify(tokenStr)
	if err != nil {
		return "", err
	}
	id, ok := claims["userId"]
	if !ok {
		return "", errors.New("can not find user id in token")
	}
	return id.(string), nil
}
