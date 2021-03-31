package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jamesoneill997/parkapi/db"
	"github.com/jamesoneill997/parkapi/initialise"
	"github.com/jamesoneill997/parkapi/logs"
	"github.com/jamesoneill997/parkapi/payments"
	"github.com/jamesoneill997/parkapi/structs"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Access struct being put here for ease of testing
type Access struct {
	l      *log.Logger
	ctx    context.Context
	client *mongo.Client
}

//NewAccess function takes a logger and constructs an access struct
func NewAccess(ctx context.Context, l *log.Logger, client *mongo.Client) *Access {
	return &Access{l, ctx, client}
}

func (user *Access) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var access structs.Access

		//read access
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logs.LogError(err)
			w.Write([]byte("Error parsing body of request"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.Unmarshal(body, &access)
		access.Uid = db.GetUserFromReg(initialise.Ctx, initialise.Client, access.Vehicle.Registration)
		currAccess := db.GetAccess(initialise.Ctx, initialise.Client, access.Uid)

		storedUser, err := db.GetUser(initialise.Ctx, initialise.Client, access.Uid.Hex())

		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error locating user"))
			return
		}

		//convert carpark id to object id
		carparkOID, err := primitive.ObjectIDFromHex(access.CarParkID)
		//user is going into the carpark
		if currAccess.Active != true {
			access.TimeStart = time.Now()
			access.Active = true

			err = db.CreateAccessInUser(initialise.Ctx, initialise.Client, access)
			if err != nil {
				logs.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error creating access"))
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Access created"))
			return
		}

		if err != nil {
			logs.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error processing carpark ID"))
			return
		}

		//if we reach this block of code, it means that the user has already accessed a carpark and hence needs to be charged
		totalCost := db.CalculateCost(initialise.Ctx, initialise.Client, currAccess.TimeStart, time.Now(), carparkOID)

		//reset access
		resetAccess := structs.Access{}
		resetAccess.Uid = access.Uid

		//update access
		err = db.CreateAccessInUser(initialise.Ctx, initialise.Client, resetAccess)

		//process charge
		payments.ChargeCustomer(storedUser["stripeid"].(string), totalCost)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("The total amount to be charged is %f", totalCost)))
		return

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}
}
