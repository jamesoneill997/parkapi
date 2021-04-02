package cors

import "net/http"

/*SetupCORS function will set the http headers required to allow communication with the front end*/
func SetupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "https://parkai.herokuapp.com/")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ParkAIToken")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}
