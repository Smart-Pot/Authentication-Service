package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	errInvalidJWT = errors.New("Invalid Google JWT")
	errInvalidIIS = errors.New("iss is invalid")
	errInvalidAUD = errors.New("aud is invalid")
	errJWTExpired =  errors.New("JWT is expired")

)	

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback",
    ClientID:     Config.Google.ClientID,
    ClientSecret: Config.Google.ClientSecret,
    Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
    Endpoint:     google.Endpoint,
}


// GoogleClaims -
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

// ValidateGoogleJWT -
func ValidateGoogleJWT(tokenString string) (GoogleClaims, error) {
	claimsStruct := GoogleClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{},  errInvalidJWT
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return GoogleClaims{}, errInvalidIIS
	}
	fmt.Println("CLAIMS",claims,"AUD",claims.Audience)
	if claims.Audience != Config.Google.ClientID {
		return GoogleClaims{}, errInvalidAUD
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return GoogleClaims{}, errJWTExpired
	}

	return *claims, nil
}


func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}