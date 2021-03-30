package db

import (
	"context"
	"parkapi/structs"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//InsertCarpark function creates a user object in the carparks collection in the database, returns mongo.InsertOne result
func InsertCarpark(ctx context.Context, client *mongo.Client, cp structs.CarPark) *mongo.InsertOneResult {

	col := client.Database("parkai").Collection("carparks")
	cp.Updated = time.Now()
	result, _ := col.InsertOne(ctx, cp)

	return result
}

/*UpdateCarpark function will update a carpark in the database to meet the spec provided*/
func UpdateCarpark(ctx context.Context, client *mongo.Client, cp structs.CarPark) bool {
	col := client.Database("parkai").Collection("carparks")

	update := bson.M{
		"$set": cp,
	}
	_, err := col.UpdateOne(
		ctx,
		bson.M{"namespace": cp.Namespace},
		update,
	)

	if err != nil {
		return false
	}

	return true
}

/*DeleteCarpark function will remove one carpark from the db, identified by namespace and UID*/
func DeleteCarpark(ctx context.Context, client *mongo.Client, namespace string, UID string) (*mongo.DeleteResult, error) {
	col := client.Database("parkai").Collection("carparks")

	//filter by owner and namespace
	filter := bson.M{
		"owner":     UID,
		"namespace": namespace,
	}

	return col.DeleteOne(ctx, filter)
}

/*DeleteCarparks function will remove a carpark from the db*/
func DeleteCarparks(ctx context.Context, client *mongo.Client, UID string) (*mongo.DeleteResult, error) {
	col := client.Database("parkai").Collection("carparks")

	//delete
	return col.DeleteMany(ctx, bson.M{"owner": UID})
}

/*FindCarParks will return a list of carparks owned by the requesting user*/
func FindCarParks(ctx context.Context, client *mongo.Client, UID string) ([]structs.CarPark, error) {
	col := client.Database("parkai").Collection("carparks")
	var carparks []structs.CarPark

	//filter by owner
	filter := bson.M{
		"owner": UID,
	}

	//run query
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return []structs.CarPark{}, err
	}

	//write all carparks to array
	if err = cursor.All(ctx, &carparks); err != nil {
		return []structs.CarPark{}, err
	}

	return carparks, nil

}

/*GetCarpark function will get a single carpark based on the user and namespace*/
func GetCarpark(ctx context.Context, client *mongo.Client, namespace string, UID string) (structs.CarPark, error) {
	col := client.Database("parkai").Collection("carparks")
	var carpark structs.CarPark

	//filter by owner and namespace
	filter := bson.M{
		"owner":     UID,
		"namespace": namespace,
	}

	//carpark does not exist
	if err := col.FindOne(ctx, filter).Decode(&carpark); err != nil {
		return structs.CarPark{}, err
	}

	return carpark, nil
}

/*CheckNamespace will check to see whether a namespace is already taken and return a boolean*/
func CheckNamespace(ctx context.Context, client *mongo.Client, namespace string, UID string) bool {
	col := client.Database("parkai").Collection("carparks")
	var carpark structs.CarPark

	//filter by owner and namespace
	filter := bson.M{
		"owner":     UID,
		"namespace": namespace,
	}

	//carpark does not exist
	if err := col.FindOne(ctx, filter).Decode(&carpark); err != nil {
		return false
	}

	//carpark exists
	return true
}

/*CalculateCost function calculates how much is due for a specific access*/
func CalculateCost(ctx context.Context, client *mongo.Client, timeStart time.Time, timeFinish time.Time, carparkID primitive.ObjectID) float32 {
	col := client.Database("parkai").Collection("carparks")
	var carpark structs.CarPark
	filter := bson.M{
		"_id": carparkID,
	}

	//write result to carpark struct
	col.FindOne(ctx, filter).Decode(&carpark)
	totalAccessTime := int(timeFinish.Sub(timeStart).Minutes())

	//minimum stay for the user in minutes
	minTime := carpark.Regulations.MinStay
	if totalAccessTime < minTime {
		return 0.0
	}

	//ceiling rounded
	chargeableHours := float32(int(totalAccessTime/60) + 1)

	//cost per hour
	costPerHour := carpark.Regulations.Cost

	return chargeableHours * costPerHour
}
