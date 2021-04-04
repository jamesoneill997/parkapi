package middleware

import (
	"errors"
	"net/http"
	"os"

	"github.com/jamesoneill997/parkapi/structs"

	"github.com/dgrijalva/jwt-go"
)

/*GetAuth will check the current actor's cookie jar to see if there is a parkai token present*/
func GetAuth(r *http.Request) (structs.Claims, error) {
	// check current session cookies
	c, err := r.Cookie("parkaitoken")

	if err != nil {
		// Unauthorised or bad request
		return structs.Claims{}, err
	}

	// Get JWT string from cookie
	tknStr := c.Value
	claims := &structs.Claims{}

	//parse error return secret env variable
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("secret")), nil
	})

	if err != nil {
		return structs.Claims{}, err

	}
	//unauthorised
	if !tkn.Valid {
		return structs.Claims{}, errors.New("InvalidToken")
	}

	return *claims, nil
}
