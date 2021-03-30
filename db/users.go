package db

import (
	"context"
	"fmt"
	"log"
	"parkapi/structs"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//InsertUser function creates a user object in the actors collection in the database, returns mongo.InsertOne result
func InsertUser(ctx context.Context, client *mongo.Client, user structs.User) *mongo.InsertOneResult {

	col := client.Database("parkai").Collection("users")

	result, insertionErr := col.InsertOne(ctx, user)

	if insertionErr != nil {
		fmt.Println("Error inserting user into database")
		log.Fatal(insertionErr)
	}

	return result
}

// FindUser function will find if a user exists in the database based on structs.user.Email
func FindUser(ctx context.Context, client *mongo.Client, email string) bool {
	var result structs.User
	col := client.Database("parkai").Collection("users")
	filter := bson.M{"email": email}

	//user does not exist
	if err := col.FindOne(ctx, filter).Decode(&result); err != nil {
		return false
	}

	//user exists
	return true
}

//GetUser will return the user that is found in the database based on structs.user.uID
func GetUser(ctx context.Context, client *mongo.Client, uID string) (bson.M, error) {
	var result bson.M
	oID, err := primitive.ObjectIDFromHex(uID)

	if err != nil {
		return result, err
	}

	col := client.Database("parkai").Collection("users")
	filter := bson.M{"_id": oID}

	//user does not exist
	if err := col.FindOne(ctx, filter).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

//UpdateUser function will process updates on a single user in the database
func UpdateUser(ctx context.Context, client *mongo.Client, user structs.User) bool {
	col := client.Database("parkai").Collection("users")

	update := bson.M{
		"$set": user,
	}

	_, err := col.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		update,
	)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

//GetType takes in a user and returns that user's type
func GetType(ctx context.Context, client *mongo.Client, user structs.User) (string, error) {
	u, err := GetUser(ctx, client, user.Email)

	if err != nil {
		return "", err
	}

	//index user.type
	actorType := u["type"]
	return fmt.Sprintf("%s", actorType), nil
}

/*DeleteUser function will remove an actor along with any vehicles or carparks that may be associated with them*/
func DeleteUser(ctx context.Context, client *mongo.Client, user structs.User) (*mongo.DeleteResult, error) {

	col := client.Database("parkai").Collection("users")
	var u structs.User

	//get type of user, get user and marshal the user
	actorType, err := GetType(ctx, client, user)
	dbUser, err := GetUser(ctx, client, user.Email)
	b, err := bson.Marshal(dbUser)

	//write the user to the user type
	bson.Unmarshal(b, &u)

	//if admin, delete all car parks from the carparks collection
	if actorType == "admin" {
		DeleteCarparks(ctx, client, user.ID.Hex())

	} else { //if user, delete all vehicles from vehicles collection
		DeleteVehicles(ctx, client, user.ID.Hex())

	}

	//delete the actor from the actors collection
	res, err := col.DeleteOne(ctx, bson.M{"_id": user.ID})
	return res, err
}

//FindUserByEmail will use an email to locate a user, can be used when user is logged out and has no claims
func FindUserByEmail(ctx context.Context, client *mongo.Client, email string) (bson.M, error) {
	var result bson.M

	col := client.Database("parkai").Collection("users")
	filter := bson.M{"email": email}

	//user does not exist
	if err := col.FindOne(ctx, filter).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

/*AddVehicleToUser adds a vehicle to users.Vehicles*/
func AddVehicleToUser(ctx context.Context, client *mongo.Client, vehicle structs.Vehicle) error {
	col := client.Database("parkai").Collection("users")

	//user that we will be adding the vehicle to
	owner, err := primitive.ObjectIDFromHex(vehicle.Owner)

	if err != nil {
		return err
	}

	//find user in users collection
	filter := bson.M{
		"_id": owner,
	}

	//update to add vehicle to vehicles.vehicles and log update on user account
	update := bson.M{
		"$push": bson.M{
			"vehicles": vehicle,
		},
		"$set": bson.M{
			"updated": time.Now(),
		},
	}

	//process update
	col.FindOneAndUpdate(ctx, filter, update)

	return nil
}

/*DeleteVehicleFromUser function deletes a vehicle from users.vehicles*/
func DeleteVehicleFromUser(ctx context.Context, client *mongo.Client, vehicle structs.Vehicle) error {
	col := client.Database("parkai").Collection("users")

	//user that we will be adding the vehicle to
	owner, err := primitive.ObjectIDFromHex(vehicle.Owner)

	if err != nil {
		return err
	}

	//find user in users collection
	filter := bson.M{
		"_id": owner,
	}

	//update to add vehicle to vehicles.vehicles and log update on user account
	update := bson.M{
		"$pull": bson.M{
			"vehicles": bson.M{
				"registration": vehicle.Registration,
			},
		},
		"$set": bson.M{
			"updated": time.Now(),
		},
	}

	//get result of update
	col.FindOneAndUpdate(ctx, filter, update)

	return nil
}

/*UpdateVehicleInUser function will update a vehicle in users.Vehicles*/
func UpdateVehicleInUser(ctx context.Context, client *mongo.Client, ud structs.RegUpdate) error {
	col := client.Database("parkai").Collection("users")

	//find user in users collection
	filter := bson.M{
		"vehicles.registration": ud.CurrReg,
	}

	update := bson.M{
		"$set": bson.M{
			"vehicles.$.registration": ud.NewReg,
			"vehicles.$.updated":      time.Now(),
			"updated":                 time.Now(),
		},
	}

	//update to add vehicle to vehicles.vehicles and log update on user account
	col.FindOneAndUpdate(ctx, filter, update)

	return nil
}

/*AddCarparkToUser adds a vehicle to users.Vehicles*/
func AddCarparkToUser(ctx context.Context, client *mongo.Client, carpark structs.CarPark) error {
	col := client.Database("parkai").Collection("users")

	//user that we will be adding the vehicle to
	owner, err := primitive.ObjectIDFromHex(carpark.Owner)

	if err != nil {
		return err
	}

	//find user in users collection
	filter := bson.M{
		"_id": owner,
	}

	//update to add vehicle to vehicles.vehicles and log update on user account
	update := bson.M{
		"$push": bson.M{
			"carparks": carpark,
		},
		"$set": bson.M{
			"updated": time.Now(),
		},
	}

	//process update
	col.FindOneAndUpdate(ctx, filter, update)

	return nil
}

/*UpdateCarparkInUser function will update a carpark in users.Carparks*/
func UpdateCarparkInUser(ctx context.Context, client *mongo.Client, cp structs.CarPark) error {
	col := client.Database("parkai").Collection("users")

	//find user in users collection
	filter := bson.M{
		"carparks.namespace": cp.Namespace,
	}

	update := bson.M{
		"$set": bson.M{
			"carparks.$.name":     cp.Name,
			"carparks.$.location": cp.Location,
			"carparks.$.rules":    cp.Regulations,
			"carparks.$.updated":  time.Now(),
			"updated":             time.Now(),
		},
	}

	//update to add carpark to carparks and log update on user account
	col.FindOneAndUpdate(ctx, filter, update)

	return nil
}

/*DeleteCarparkFromUser function deletes a carpark from users.carparks*/
func DeleteCarparkFromUser(ctx context.Context, client *mongo.Client, carpark structs.CarPark, uid string) error {
	col := client.Database("parkai").Collection("users")

	//user that we will be removing the carpark from
	owner, err := primitive.ObjectIDFromHex(uid)

	if err != nil {
		return err
	}

	//find user in users collection
	filter := bson.M{
		"_id": owner,
	}

	//update to remove carpark and update user
	update := bson.M{
		"$pull": bson.M{
			"carparks": bson.M{
				"namespace": carpark.Namespace,
			},
		},
		"$set": bson.M{
			"updated": time.Now(),
		},
	}

	//get result of update
	col.FindOneAndUpdate(ctx, filter, update)

	return nil
}

/*GetAccess function will check whether the current user is currently accessing a carpark*/
func GetAccess(ctx context.Context, client *mongo.Client, uid primitive.ObjectID) structs.Access {
	col := client.Database("parkai").Collection("users")
	var user structs.User

	//locate user by id
	filter := bson.M{
		"_id": uid,
	}

	//write result to user struct
	col.FindOne(ctx, filter).Decode(&user)

	return user.Access
}

/*CreateAccessInUser function will create an access object in users.Access*/
func CreateAccessInUser(ctx context.Context, client *mongo.Client, access structs.Access) error {
	col := client.Database("parkai").Collection("users")
	//locate user by id
	filter := bson.M{
		"_id": access.Uid,
	}

	//update to add access to user.Access
	update := bson.M{
		"$set": bson.M{
			"access":  access,
			"updated": time.Now(),
		},
	}

	//process update
	result := col.FindOneAndUpdate(ctx, filter, update)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

/*GetUserFromReg will get a user from a registration plate number that is passed to the function*/
func GetUserFromReg(ctx context.Context, client *mongo.Client, reg string) primitive.ObjectID {
	col := client.Database("parkai").Collection("users")
	var owner structs.User
	filter := bson.M{
		"vehicles.registration": reg,
	}

	col.FindOne(ctx, filter).Decode(&owner)

	return owner.ID
}
