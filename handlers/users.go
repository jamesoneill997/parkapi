package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jamesoneill997/parkapi/cors"
	"github.com/jamesoneill997/parkapi/db"
	"github.com/jamesoneill997/parkapi/initialise"
	"github.com/jamesoneill997/parkapi/logs"
	"github.com/jamesoneill997/parkapi/middleware"
	"github.com/jamesoneill997/parkapi/payments"
	"github.com/jamesoneill997/parkapi/structs"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User struct being put here for ease of testing
type User struct {
	l *log.Logger
}

//NewUser function takes a logger and constructs a user struct with that logger as User.l
func NewUser(l *log.Logger) *User {
	return &User{l}
}

//UsersHandler function will handle calls to the users endpoint
func (user *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//allow CORS
	cors.SetupCORS(&w, r)

	var actor structs.User
	switch r.Method {

	//handling CORS
	case http.MethodOptions:
		w.WriteHeader(200)
		return

	//handles get requests
	case http.MethodGet:
		//get authorisation
		claims, err := middleware.GetAuth(r)

		//handle auth err
		if err != nil {
			if err == jwt.ErrSignatureInvalid || err == http.ErrNoCookie || err.Error() == "InvalidToken" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorised"))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		dbUser, err := db.GetUser(initialise.Ctx, initialise.Client, claims.Uid)
		bsonUser, err := bson.Marshal(dbUser)

		bson.Unmarshal(bsonUser, &actor)
		jsonUser, err := json.Marshal(actor)

		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error processing request"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s", jsonUser)))

		return

	//handles post request
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logs.LogError(err)
			w.Write([]byte("Error parsing body of request"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.Unmarshal(body, &actor)

		//log current time and use as created/updated values
		actor.Created = time.Now()
		actor.Updated = time.Now()

		//hash password before storage
		if actor.Password, err = middleware.HashPass(actor.Password); err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error processing your request, please try again later."))
			return
		}

		if !db.FindUser(initialise.Ctx, initialise.Client, actor.Email) {
			//create and add id to actor
			id := primitive.NewObjectID()
			actor.ID = id

			//add vehicle or carpark to respective collections
			if actor.Type == "user" {
				vehicles := actor.Vehicles
				for i, v := range vehicles {
					//generate random unique id and assign for each vehicle
					id = primitive.NewObjectID()
					v.ID = id
					v.Owner = actor.ID.Hex()

					//log creation and update times
					v.Created = time.Now()
					v.Updated = time.Now()

					//set actor vehicle to copy of vehicle after above updates
					actor.Vehicles[i] = v

					//insert to vehicles collection
					db.InsertVehicle(initialise.Ctx, initialise.Client, v)
				}

			} else {
				cp := actor.CarParks
				for i, c := range cp {
					//generate random unique id and assign for each vehicle
					id = primitive.NewObjectID()
					c.ID = id
					c.Owner = actor.ID.Hex()

					//log creation and update times
					c.Created = time.Now()
					c.Updated = time.Now()

					actor.CarParks[i] = c

					//insert to carparks collection
					db.InsertCarpark(initialise.Ctx, initialise.Client, c)
				}
			}
			customer, err := payments.CreateCustomer(actor)

			if err != nil {
				logs.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error communicating with payments service"))
				return
			}

			//add customer id to locally stored object
			actor.StripeID = customer.ID

			db.InsertUser(initialise.Ctx, initialise.Client, actor)

			//create json web token for auth
			jwt, err := middleware.CreateJWT(actor.ID.Hex())

			if err != nil {
				logs.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error processing your request, please try again later."))
				return

			}

			//set cookie, valid for 15 minutes
			middleware.SetCookie("ParkAIToken", jwt, time.Now().Add(15*time.Minute), w)

			//send welcome email to user
			middleware.WelcomeEmail(actor)
			jsonUser, err := json.Marshal(actor)

			//success
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(fmt.Sprintf("%s", jsonUser)))
			return

		}
		//User already exists
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("User with that email address already exists."))
		return

	//patch will handle updates to an actor
	case http.MethodPatch:
		body, err := ioutil.ReadAll(r.Body)

		//internal error, send to log file
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing body of request"))
			return
		}

		//initialise 2 users to be compared
		var actor structs.User
		var storedUser *structs.User

		//get authorisation
		claims, err := middleware.GetAuth(r)

		//handle auth err
		if err != nil {
			if err == jwt.ErrSignatureInvalid || err == http.ErrNoCookie || err.Error() == "InvalidToken" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorised"))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		//get db entry for currently logged in user
		dbUser, err := db.GetUser(initialise.Ctx, initialise.Client, claims.Uid)
		bsonUser, _ := bson.Marshal(dbUser)

		//stored actor and request body actor
		bson.Unmarshal(bsonUser, &storedUser)
		json.Unmarshal(body, &actor)

		//check that updated email is not already in use
		if db.FindUser(initialise.Ctx, initialise.Client, actor.Email) && actor.Email != storedUser.Email {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Email already in use"))
			return
		}

		//log current time and use as created/updated values
		actor.Updated = time.Now()

		//created time doesn't change
		actor.Created = storedUser.Created

		//update fields on actor object to be consistent with database and secure
		actor.Password, err = middleware.HashPass(actor.Password)
		aid, err := primitive.ObjectIDFromHex(claims.Uid)
		actor.ID = aid
		actor.StripeID = storedUser.StripeID

		//internal error, send to log file
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error updating account. Please try again later."))
			return
		}

		//boolean variable to check whether the update was successful or not
		var updateSuccessful bool

		//if actor, then update carparks in carparks collection
		if actor.Type == "admin" {
			for i := range actor.CarParks {
				actor.CarParks[i].Owner = claims.Uid
				actor.CarParks[i].ID = storedUser.CarParks[i].ID
				actor.CarParks[i].Created = storedUser.CarParks[i].Created
				actor.CarParks[i].Updated = time.Now()

				updateSuccessful = db.UpdateCarpark(initialise.Ctx, initialise.Client, actor.CarParks[i])
			}
		} else { //if we are dealing with a user, then we update the vehicles collection in the database
			for i := range actor.Vehicles {
				actor.Vehicles[i].Owner = claims.Uid
				actor.Vehicles[i].ID = storedUser.Vehicles[i].ID
				actor.Vehicles[i].Created = storedUser.Vehicles[i].Created
				actor.Vehicles[i].Updated = time.Now()

				updateSuccessful = db.UpdateVehicle(initialise.Ctx, initialise.Client, actor.Vehicles[i])
			}
		}

		//if there was an error processing the update to either database
		if !updateSuccessful {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error updating account. Please try again later."))
			return
		}

		//now update users collection
		updateSuccessful = db.UpdateUser(initialise.Ctx, initialise.Client, actor)

		//check if there was an issue updating the user
		if !updateSuccessful {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error updating account. Please try again later."))
			return
		}

		//no content, update successful
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Resource updated successfully"))

	//delete request deletes user data in all collections, including carparks and vehicles
	case http.MethodDelete:
		var u structs.User

		//get authorisation
		claims, err := middleware.GetAuth(r)

		//handle auth err
		if err != nil {
			if err == jwt.ErrSignatureInvalid || err == http.ErrNoCookie || err.Error() == "InvalidToken" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorised"))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		//setup user to match data in db
		currUserID, err := primitive.ObjectIDFromHex(claims.Uid)
		u.ID = currUserID
		u.Email = claims.Email

		//error here means we cannot identify user
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing user identity from token"))
			return
		}

		storedUser, err := db.GetUser(initialise.Ctx, initialise.Client, claims.Uid)
		_, err = payments.DeleteCustomer(storedUser["stripeid"].(string))
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error removing user from payments service"))
			return
		}
		//process delete request
		_, err = db.DeleteUser(initialise.Ctx, initialise.Client, u)

		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error removing user from system"))
			return
		}

		//error here means that there is an issue with the db delete function
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Issue processing delete request, please try again later"))
			return
		}

		//request valid and complete
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Resource successfully deleted"))

	//all other requests ignored
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
