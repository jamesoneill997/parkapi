package handlers

import (
	"context"
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
	"github.com/jamesoneill997/parkapi/structs"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Vehicle struct being put here for ease of testing
type Vehicle struct {
	l      *log.Logger
	ctx    context.Context
	client *mongo.Client
}

//NewVehicle function takes a logger and constructs a carpark struct
func NewVehicle(ctx context.Context, l *log.Logger, client *mongo.Client) *Vehicle {
	return &Vehicle{l, ctx, client}
}

func (user *Vehicle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//allow CORS
	cors.SetupCORS(&w, r)

	switch r.Method {
	case http.MethodGet:
		claims, err := middleware.GetAuth(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised"))
			return
		}

		vehicles, err := db.GetVehicles(initialise.Ctx, initialise.Client, claims.Uid)
		jsonUser, err := json.Marshal(vehicles)

		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error processing request. Please try again later"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s", jsonUser)))
		return

	case http.MethodPost:
		claims, err := middleware.GetAuth(r)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised"))
			return
		}

		var v structs.Vehicle
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing body of request"))
			return
		}

		json.Unmarshal(body, &v)

		//check if the vehicle has already been registered on the db
		if db.CheckRegistration(initialise.Ctx, initialise.Client, v.Registration) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Vehicle with this registration number already exists"))
			return
		}

		//set owner to current user
		v.Owner = claims.Uid

		//generate random unique id and assign for each vehicle
		id := primitive.NewObjectID()
		v.ID = id

		//log creation and update times
		v.Created = time.Now()
		v.Updated = time.Now()

		//insert vehicle to vehicles collection, then users collection
		db.InsertVehicle(initialise.Ctx, initialise.Client, v)
		err = db.AddVehicleToUser(initialise.Ctx, initialise.Client, v)

		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error adding vehicle, please try again later"))
			return
		}

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
		var details structs.RegUpdate
		body, err := ioutil.ReadAll(r.Body)

		//internal error
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing body of request"))
			return
		}

		//write body to details
		json.Unmarshal(body, &details)

		//check that vehicle exists and belongs to requesting user
		vehicle, err := db.GetVehicle(initialise.Ctx, initialise.Client, details.CurrReg)
		if err != nil || claims.Uid != vehicle.Owner {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unable to locate resource"))
			return
		}

		//update reg
		vehicle.Registration = details.NewReg
		if db.CheckRegistration(initialise.Ctx, initialise.Client, details.NewReg) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Vehicle with this registration number already exists"))
			return
		}

		//update vehicle in vehicles and users collection
		db.UpdateVehicle(initialise.Ctx, initialise.Client, vehicle)
		err = db.UpdateVehicleInUser(initialise.Ctx, initialise.Client, details)

		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error updating vehicle"))
			return
		}

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
		var details structs.Vehicle
		body, err := ioutil.ReadAll(r.Body)

		//internal error
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing body of request"))
			return
		}

		//write body to details
		json.Unmarshal(body, &details)

		//check that vehicle exists and belongs to requesting user
		vehicle, err := db.GetVehicle(initialise.Ctx, initialise.Client, details.Registration)
		if err != nil || claims.Uid != vehicle.Owner {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unable to locate resource"))
			return
		}

		//remove vehicle from vehicles collection and users collection
		db.DeleteVehicle(initialise.Ctx, initialise.Client, vehicle.Registration)
		err = db.DeleteVehicleFromUser(initialise.Ctx, initialise.Client, vehicle)
		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error deleting vehicle"))
			return
		}
		//request successful
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Resource Deleted Successfully"))
		return

	//reject all other requests
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}
}
