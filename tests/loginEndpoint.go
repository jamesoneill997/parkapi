package testing

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jamesoneill997/parkapi/handlers"
	"github.com/jamesoneill997/parkapi/initialise"
	"github.com/jamesoneill997/parkapi/structs"

	"github.com/stretchr/testify/assert"
)

/*Test suite to run testing on login endpoint*/

/*LoginEndpointTest runs all test files on the login endpoint*/
func LoginEndpointTest(t *testing.T, users structs.Users) {
	//creds type that can be written to from user struct
	var creds structs.LoginCreds

	//range over users and perform checks on each one
	for _, u := range users.Users {
		creds = GetCreds(creds, u)
		loginValid(t, creds)
		loginInvalidEmail(t, creds)
		loginInvalidPassword(t, creds)
	}

}

/*GetCreds function gets credentials from test users and stores them in a LoginCreds struct*/
func GetCreds(creds structs.LoginCreds, u structs.User) structs.LoginCreds {
	//create json object and write it to Login Creds, this extracts the necessary fields from the user object
	jsonUser, _ := json.Marshal(u)
	json.Unmarshal(jsonUser, &creds)

	return creds
}

//loginValid function will attempt to login a user with correct credentials
func loginValid(t *testing.T, loginCreds structs.LoginCreds) {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//create json user for POST request
	jsonUser, err := json.Marshal(loginCreds)
	if err != nil {
		return
	}

	//send request
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonUser))
	response := httptest.NewRecorder()
	handlers.NewLogin(initialise.Ctx, l, initialise.Client).ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response expected")
}

//loginInvalidEmail function will attempt to login with an invalid email
func loginInvalidEmail(t *testing.T, loginCreds structs.LoginCreds) {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//manipulate email so it's incorrect
	loginCreds.UserEmail += "invalidStringToAppend"

	//create json user for POST request
	jsonUser, err := json.Marshal(loginCreds)
	if err != nil {
		return
	}

	//send request
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonUser))
	response := httptest.NewRecorder()
	handlers.NewLogin(initialise.Ctx, l, initialise.Client).ServeHTTP(response, request)
	assert.Equal(t, 401, response.Code, "Unauthorised response expected")
}

//loginInvalidPassword will attempt to login with an invalid password
func loginInvalidPassword(t *testing.T, loginCreds structs.LoginCreds) {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//manipulate the password so it's incorrect
	loginCreds.UserPw += "invalidStringToAppend"

	//create json user for POST request
	jsonUser, err := json.Marshal(loginCreds)
	if err != nil {
		return
	}

	//send request
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonUser))
	response := httptest.NewRecorder()
	handlers.NewLogin(initialise.Ctx, l, initialise.Client).ServeHTTP(response, request)
	assert.Equal(t, 401, response.Code, "Unauthorised response expected")
}
