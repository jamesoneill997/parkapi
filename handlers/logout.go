package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jamesoneill997/parkapi/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

//Logout struct being put here for ease of testing
type Logout struct {
	l      *log.Logger
	ctx    context.Context
	client *mongo.Client
}

//NewLogout function takes a logger and constructs a login struct
func NewLogout(ctx context.Context, l *log.Logger, client *mongo.Client) *Logout {
	return &Logout{l, ctx, client}
}

func (user *Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//allow CORS
	cors.SetupCORS(&w, r)

	switch r.Method {
	//api should only accept POST request for logout
	case http.MethodPost:
		//variables to create cookie
		name := "ParkAIToken"
		token := ""
		expiration := time.Now()

		//set cookie to expire now
		http.SetCookie(w, &http.Cookie{
			Name:    name,
			Value:   token,
			Expires: expiration,
		},
		)

		//no content, request valid and served
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Successfully logged out"))

	//all other requests should be rejected
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}
}
