package cors

import "net/http"

/*SetupCORS function will set the http headers required to allow communication with the front end*/
func SetupCORS(w *http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")

	(*w).Header().Set("Access-Control-Allow-Origin", origin)
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Set-Cookie, set-cookie")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}
