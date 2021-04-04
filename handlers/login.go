package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jamesoneill997/parkapi/cors"
	"github.com/jamesoneill997/parkapi/db"
	"github.com/jamesoneill997/parkapi/logs"
	"github.com/jamesoneill997/parkapi/middleware"
	"github.com/jamesoneill997/parkapi/structs"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Login struct being put here for ease of testing
type Login struct {
	l      *log.Logger
	ctx    context.Context
	client *mongo.Client
}

//NewLogin function takes a logger and constructs a login struct
func NewLogin(ctx context.Context, l *log.Logger, client *mongo.Client) *Login {
	return &Login{l, ctx, client}
}

func (user *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//struct to store login credentials
	var details structs.LoginCreds

	//allow CORS
	cors.SetupCORS(&w, r)

	//handle requests
	switch r.Method {

	//handling CORS
	case http.MethodOptions:
		w.WriteHeader(200)
		return

	//post is the only valid request for this endpoint
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logs.LogError(err)
			w.Write([]byte("Error parsing body of request"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//parse input to json
		e := json.Unmarshal(body, &details)
		if e != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing body of request"))
			return
		}

		//get user from database
		user, err := db.FindUserByEmail(user.ctx, user.client, details.UserEmail)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("The credentials entered were incorrect"))
			return
		}
		pw := user["password"]

		//verify password on db with password entered
		authenticated := middleware.CheckPassword(details.UserPw, pw.(string))
		if !authenticated {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Credentials are incorrect"))
			return
		}

		uid := user["_id"]
		//create new key
		token, err := middleware.CreateJWT(uid.(primitive.ObjectID).Hex())
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("An error has occurred with your request. Please try again later."))
		}

		//set cookie
		middleware.SetCookie("ParkAIToken", token, time.Now().Add(48*time.Hour), w)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
	}
}
