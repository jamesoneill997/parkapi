package middleware

import (
	"net/http"
	"os"
	"parkapi/structs"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte(os.Getenv("secret"))

//CreateJWT will create a json web token from the unique ID provided in the user's profile
func CreateJWT(ID string) (string, error) {
	expire := time.Now().Add(15 * time.Minute)
	claims := &structs.Claims{
		Uid: ID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expire.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

//SetCookie will set the cookie for the end user
func SetCookie(name string, token string, expiration time.Time, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   token,
		Expires: expiration,
	})
}
