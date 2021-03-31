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
	"github.com/jamesoneill997/parkapi/structs"

	"github.com/stretchr/testify/assert"
)

/*This is the test suite for the user endpoint of the api, it handles tests on the CRUD operations of actors*/

/*UserEndpointTests runs all of the tests below*/
func UserEndpointTests(t *testing.T, testUsers structs.Users) {
	getLoggedOut(t)

	//test signup for all of our test users
	for i := 0; i < len(testUsers.Users); i++ {
		signUp(t, testUsers.Users[i])
		duplicateSignUp(t, testUsers.Users[i])
	}
}

/*Test functions*/
//getLoggedOut tests the case where the user is logged out (no token present) and sends a get request.
func getLoggedOut(t *testing.T) {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//send GET request to users endpoint
	request, _ := http.NewRequest("GET", "/users", nil)
	response := httptest.NewRecorder()
	handlers.NewUser(l).ServeHTTP(response, request)

	//as user does not have a token, they are not authorised
	assert.Equal(t, 401, response.Code, "Unauthorised response expected")
}

//signUp tests the response of a valid sign up request
func signUp(t *testing.T, user structs.User) {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//create json user for request
	jsonUser, _ := json.Marshal(user)

	//send POST request with json user
	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonUser))
	response := httptest.NewRecorder()
	handlers.NewUser(l).ServeHTTP(response, request)

	//Resource should be created
	assert.Equal(t, 201, response.Code, "Created response expected")
}

//duplicate sign up is identical to sign up, but is always run after signUp and expects a 409
func duplicateSignUp(t *testing.T, user structs.User) {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//send POST request to users endpoint
	jsonUser, _ := json.Marshal(user)
	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonUser))
	response := httptest.NewRecorder()
	handlers.NewUser(l).ServeHTTP(response, request)

	//user should already exist
	assert.Equal(t, 409, response.Code, "Resource already exists result expected")
}
