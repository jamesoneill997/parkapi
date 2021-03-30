package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"parkapi/db"
	"parkapi/initialise"
	"parkapi/logs"
	"parkapi/middleware"
	"parkapi/structs"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//CarPark struct being put here for ease of testing
type CarPark struct {
	l      *log.Logger
	ctx    context.Context
	client *mongo.Client
}

//NewCarPark function takes a logger and constructs a carpark struct
func NewCarPark(ctx context.Context, l *log.Logger, client *mongo.Client) *CarPark {
	return &CarPark{l, ctx, client}
}

func (user *CarPark) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		claims, err := middleware.GetAuth(r)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised"))
			return
		}

		// find users carparks in the database
		cps, err := db.FindCarParks(initialise.Ctx, initialise.Client, claims.Uid)
		jsonCPs, err := json.Marshal(cps)

		//internal error
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error processing request. Please try again later"))
			return
		}

		//tell receiver the content type will be JSON
		w.Header().Set("Content-Type", "application/json")

		//Status ok
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s", jsonCPs)))
		return

	case http.MethodPost:
		//check user auth
		claims, err := middleware.GetAuth(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised"))
			return
		}

		//init carpark
		var cp structs.CarPark

		//read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing body of request"))
			return
		}

		//write body to carpark
		json.Unmarshal(body, &cp)

		//set undefined fields
		cp.Owner = claims.Uid
		id := primitive.NewObjectID()
		cp.ID = id

		//log creation and update times
		cp.Created = time.Now()
		cp.Updated = time.Now()

		//add cp to carparks and users collections
		db.InsertCarpark(initialise.Ctx, initialise.Client, cp)
		err = db.AddCarparkToUser(initialise.Ctx, initialise.Client, cp)

		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error inserting carpark."))
		}

		//status created
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Resource Created Successfully"))
		return

	case http.MethodPatch:
		//get auth for user
		claims, err := middleware.GetAuth(r)

		//handle unauthorised request
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised"))
			return
		}

		//initialise struct and read in request body
		var cp structs.CarPark
		body, err := ioutil.ReadAll(r.Body)

		//internal error
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing body of request"))
			return
		}

		//write body to details
		json.Unmarshal(body, &cp)

		//check that vehicle exists and belongs to requesting user
		storedCp, err := db.GetCarpark(initialise.Ctx, initialise.Client, cp.Namespace, claims.Uid)
		if err != nil || claims.Uid != storedCp.Owner {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unable to locate resource"))
			return
		}

		//update unspecified fields
		cp.Updated = time.Now()
		cp.ID = storedCp.ID
		cp.Owner = storedCp.Owner
		cp.Created = storedCp.Created

		//process update
		db.UpdateCarpark(initialise.Ctx, initialise.Client, cp)
		db.UpdateCarparkInUser(initialise.Ctx, initialise.Client, cp)

		//no content, update successful
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Resource Updated Successfully"))
		return

	case http.MethodDelete:
		//get auth for user
		claims, err := middleware.GetAuth(r)

		//handle unauthorised request
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised"))
			return
		}

		//initialise struct and read in request body
		var carpark structs.CarPark
		body, err := ioutil.ReadAll(r.Body)

		//internal error
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing body of request"))
			return
		}

		//write body to details
		json.Unmarshal(body, &carpark)

		//check that carpark exists and belongs to requesting user
		storedCp, err := db.GetCarpark(initialise.Ctx, initialise.Client, carpark.Namespace, claims.Uid)

		if err != nil || claims.Uid != storedCp.Owner {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unable to locate resource"))
			return
		}

		//remove carpark from db
		db.DeleteCarpark(initialise.Ctx, initialise.Client, carpark.Namespace, claims.Uid)
		db.DeleteCarparkFromUser(initialise.Ctx, initialise.Client, carpark, claims.Uid)

		//Status Ok
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Resource Deleted Successfully"))
		return

	//All other requests should be ignored
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}
}
