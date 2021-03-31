package testing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/jamesoneill997/parkapi/structs"
)

var r *http.Request

//TestMain will test the user endpoint of the api with multiple scenarios
func TestMain(t *testing.T) {
	//get test users from testUsers.json
	testUsers := getTestUsers()

	/*Tests*/
	//UserEndpointTests(t, testUsers)
	LoginEndpointTest(t, testUsers)
	LogoutEndpointTest(t)
}

//GetTestUsers reads in test user data from json file
func getTestUsers() structs.Users {
	//struct to store data
	var result structs.Users

	// Open file
	jsonFile, err := os.Open("./testUsers.json")
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)

	fmt.Println("Successfully fetched test users from file")

	return result
}
