package db

import (
	"context"
	"fmt"
	"parkapi/logs"
	"parkapi/structs"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//InsertVehicle function creates a user object in the vehicles collection in the database, returns mongo.InsertOne result
func InsertVehicle(ctx context.Context, client *mongo.Client, v structs.Vehicle) *mongo.InsertOneResult {
	col := client.Database("parkai").Collection("vehicles")
	result, insertionErr := col.InsertOne(ctx, v)

	if insertionErr != nil {
		fmt.Println("Error inserting vehicle into database")
		logs.LogError(insertionErr)
	}

	return result
}

//UpdateVehicle function will process updates on a single user in the database
func UpdateVehicle(ctx context.Context, client *mongo.Client, v structs.Vehicle) bool {
	col := client.Database("parkai").Collection("vehicles")

	update := bson.M{
		"$set": bson.M{
			"registration": v.Registration,
			"updated":      time.Now(),
		},
	}
	_, err := col.UpdateOne(
		ctx,
		bson.M{"_id": v.ID},
		update,
	)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

/*DeleteVehicles function will remove vehicles from the db*/
func DeleteVehicles(ctx context.Context, client *mongo.Client, UID string) (*mongo.DeleteResult, error) {
	col := client.Database("parkai").Collection("vehicles")

	//delete
	return col.DeleteMany(ctx, bson.M{"owner": UID})
}

/*GetVehicles will return all vehicles that belong to the requesting user*/
func GetVehicles(ctx context.Context, client *mongo.Client, UID string) ([]structs.Vehicle, error) {
	var vehicles []structs.Vehicle

	col := client.Database("parkai").Collection("vehicles")

	//filter by owner
	filter := bson.M{
		"owner": UID,
	}

	cursor, err := col.Find(ctx, filter)

	if err != nil {
		return []structs.Vehicle{}, err
	}

	if err = cursor.All(ctx, &vehicles); err != nil {
		return []structs.Vehicle{}, err
	}

	return vehicles, nil
}

/*CheckRegistration will check whether the requested reg plate is already in the db or not, if it is then we need to reject the request*/
func CheckRegistration(ctx context.Context, client *mongo.Client, reg string) bool {
	var result structs.Vehicle
	col := client.Database("parkai").Collection("vehicles")
	filter := bson.M{"registration": reg}

	//reg does not exist
	if err := col.FindOne(ctx, filter).Decode(&result); err != nil {
		return false
	}

	//reg exists
	return true
}

/*GetVehicle searches the database for a vehicle based on the reg and returns an object of type structs.Vehicle*/
func GetVehicle(ctx context.Context, client *mongo.Client, reg string) (structs.Vehicle, error) {
	var result structs.Vehicle

	col := client.Database("parkai").Collection("vehicles")
	filter := bson.M{
		"registration": reg,
	}

	//reg does not exist
	if err := col.FindOne(ctx, filter).Decode(&result); err != nil {
		return structs.Vehicle{}, err
	}

	return result, nil
}

/*DeleteVehicle function will delete a vehicle keyed by reg*/
func DeleteVehicle(ctx context.Context, client *mongo.Client, reg string) (*mongo.DeleteResult, error) {
	col := client.Database("parkai").Collection("vehicles")
	filter := bson.M{"registration": reg}

	return col.DeleteOne(ctx, filter)
}
