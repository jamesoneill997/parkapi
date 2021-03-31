package testing

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jamesoneill997/parkapi/handlers"
	"github.com/jamesoneill997/parkapi/initialise"

	"github.com/stretchr/testify/assert"
)

//Test suite for logout functionality

/*LogoutEndpointTest will test the logout endpoint*/
func LogoutEndpointTest(t *testing.T) {
	logout(t)
	logoutInvalid(t)
}

//logout request, should give same status code whether user is logged in or not
func logout(t *testing.T) {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//send a POST request to logout endpoint
	request, _ := http.NewRequest("POST", "/logout", nil)
	response := httptest.NewRecorder()
	handlers.NewLogout(initialise.Ctx, l, initialise.Client).ServeHTTP(response, request)
	assert.Equal(t, http.StatusNoContent, response.Code, "No content response expected")
}

//logout GET request, should return bad request
func logoutInvalid(t *testing.T) {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//Send GET request to logout enpoint
	request, _ := http.NewRequest("GET", "/logout", nil)
	response := httptest.NewRecorder()
	handlers.NewLogout(initialise.Ctx, l, initialise.Client).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Bad Request response expected")
}
