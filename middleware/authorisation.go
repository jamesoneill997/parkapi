package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jamesoneill997/parkapi/structs"

	"github.com/dgrijalva/jwt-go"
)

/*GetAuth will check the current actor's cookie jar to see if there is a parkai token present*/
func GetAuth(r *http.Request) (structs.Claims, error) {
	var tknStr string
	var claims = structs.Claims{}
	c := jwt.MapClaims{}

	// check current session cookies
	_, err := r.Cookie("ParkAIToken")

	if err != nil {
		// Get JWT string from cookie
		tknStrArr := strings.Split(r.Header.Get("Set-Cookie"), "=")
		if len(tknStrArr) > 1 {
			tknStr = tknStrArr[1]
		} else {
			return claims, err
		}
	}

	//parse error return secret env variable
	tkn, err := jwt.ParseWithClaims(tknStr, c, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("secret")), nil
	})

	if err != nil {
		fmt.Println(err)
		return claims, err

	}
	//unauthorised
	if !tkn.Valid {
		return structs.Claims{}, errors.New("InvalidToken")
	}

	return claims, nil
}
