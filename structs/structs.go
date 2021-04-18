package structs

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
User struct stores user details
*/
type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	Type           string             `json:"type" bson: "type"`
	FirstName      string             `json:"firstName" bson:"firstName"`
	Surname        string             `json:"surname" bson:"surname"`
	Email          string             `json:"email" bson:"email"`
	Password       string             `json:"password" bson:"password"`
	StripeID       string             `json:"stripeid" bson:"stripeid"`
	AccountBalance float32            `json:"accountBalance" bson:"accountBalance"`
	CarParks       []CarPark          `json:"carparks" bson:"carparks"`
	Vehicles       []Vehicle          `json:"vehicles" bson:"vehicles"`
	Access         Access             `json: "access" bson:"access"`
	Created        time.Time          `json:"created" bson:"created"`
	Updated        time.Time          `json:"updated" bson:"updated"`
}

/*
Users stores list of users
*/
type Users struct {
	Users []User `json:"users" bson:"users"`
}

/*
Vehicle struct will store vehicle details
*/
type Vehicle struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Registration string             `json:"registration" bson:"registration"`
	Owner        string             `json:"owner" bson:"owner"` //user ID
	Created      time.Time          `json:"created" bson:"created"`
	Updated      time.Time          `json:"updated" bson:"updated"`
}

/*LoginCreds stores the login credentials that are passed to the api during a login request*/
type LoginCreds struct {
	UserType  string `json:"type"`
	UserEmail string `json:"email"`
	UserPw    string `json:"password"`
}

/*
CarPark struct stores CarPark details
*/
type CarPark struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Namespace   string             `json:"namespace" bson:"namespace"` //unique identifier for carpark
	Location    Address            `json:"location" bson:"location"`
	Regulations Settings           `json:"rules" bson:"rules"`
	Owner       string             `json:"owner" bson:"owner"` //Admin ID
	Created     time.Time          `json:"created" bson:"created"`
	Updated     time.Time          `json:"updated" bson:"updated"`
}

/*
Settings struct stores CarPark settings details. This is only referenced in the CarPark struct
*/
type Settings struct {
	Hours        [7][2]int `json:"hours" bson:"hours"` //one array for each day of the week, nested array first element is opening time, last element is closing time (int, military time)
	Capacity     int       `json:"capacity" bson:"capacity"`
	BlockedUsers []string  `json:"blockedUsers" bson:"blockedUsers"` //list of user ID's
	MinStay      int       `json:"minStay" bson:"minStay"`           //minimum stay in minutes
	MaxStay      int       `json:"maxStay" bson:"maxStay"`           //maximum stay in hours
	Cost         float32   `json:"cost" bson:"cost"`                 //cost per hour
}

/*
Address struct stores address details which are used in CarPark
*/
type Address struct {
	Street   string `json:"street" bson:"street"`
	Town     string `json:"town" bson:"town"`
	County   string `json:"county" bson:"county"`
	Postcode string `json:"postcode" bson:"postcode"`
}

//Claims for authentication purposes (JWT)
type Claims struct {
	Uid            string    `json: "uid" bson: "uid"`
	Email          string    `json:"email" bson:"email"`
	Type           string    `json:"type" bson: "type"`
	FirstName      string    `json:"firstName" bson:"firstName"`
	Surname        string    `json:"surname" bson:"surname"`
	Password       string    `json:"password" bson:"password"`
	StripeID       string    `json:"stripeid" bson:"stripeid"`
	AccountBalance float32   `json:"accountBalance" bson:"accountBalance"`
	CarParks       []CarPark `json:"carparks" bson:"carparks"`
	Vehicles       []Vehicle `json:"vehicles" bson:"vehicles"`
	Access         Access    `json: "access" bson:"access"`
	Created        time.Time `json:"created" bson:"created"`
	Updated        time.Time `json:"updated" bson:"updated"`
	jwt.StandardClaims
}

/*Access struct is used to store data relating to a carpark access*/
type Access struct {
	Active    bool               `json: "active" bson:"active"` //false by default, true if user is currently in carpark
	CarParkID string             `json: "carparkid" bson: "carparkid"`
	Uid       primitive.ObjectID `json:"uid" bson: "uid"`
	Vehicle   Vehicle            `json:"vehicle" bson:"vehicle"`
	TimeStart time.Time          `json:"start" bson:"start"`
}

/*RegUpdate struct stores details for the update of vehicle.Registration*/
type RegUpdate struct {
	CurrReg string `json: "currreg" bson: "curreg"`
	NewReg  string `json: "newreg" bson: "newreg"`
}
