package initialise

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Store Client and Ctx of database connection for external ref
var Client, Ctx = Connection(os.Getenv("mongosrv"))

//Connection function connects to the mongo cluster using the URI. URI passed as parameter
func Connection(uri string) (*mongo.Client, context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Error connecting to MongoDB")
	}

	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Error connecting to MongoDB")
	}

	//check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Error connecting to MongoDB")
	}

	//check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Cannot connect to MongoDB")
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println("Connection to Mongodb timed out")
		log.Fatal(err)
	}

	fmt.Println("Connected to Mongodb")

	return client, ctx
}
