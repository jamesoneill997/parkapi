package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jamesoneill997/parkapi/db"
	"github.com/jamesoneill997/parkapi/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte(os.Getenv("secret"))

//CreateJWT will create a json web token from the unique ID provided in the user's profile
func CreateJWT(ID string, ctx context.Context, client *mongo.Client) (string, error) {
	expire := time.Now().Add(15 * time.Minute)
	var actor structs.User

	dbUser, err := db.GetUser(ctx, client, ID)
	bsonUser, err := bson.Marshal(dbUser)

	bson.Unmarshal(bsonUser, &actor)
	jsonUser, err := json.Marshal(actor)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonUser)

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
		Path:    "/",
		Expires: expiration,
	})
}
