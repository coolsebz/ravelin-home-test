package main

import (
	"github.com/coolsebz/ravelin-home-test/backend/handlers"
	"log"
	"net/http"
)

// type for describing more than just endpoints
// putting together endpoint + HTTP method
type RestResource struct {
	Route  string
	Method string
}

// a representation of all the handlers we're exposing mapped to routes
var mux map[RestResource]func(http.ResponseWriter, *http.Request)

type customHandler struct{}

func (*customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// doing our CORS setup prior to reaching the handlers
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	if handler, ok := mux[RestResource{Route: r.URL.String(), Method: r.Method}]; ok {
		handler(w, r)
		return
	}

}

// creating and adding our route/handler mappings
func setupHandlers() {
	mux = make(map[RestResource]func(http.ResponseWriter, *http.Request))

	// our handler for when a new "event" is sent to the server
	mux[RestResource{Route: "/events", Method: http.MethodPost}] = handlers.ReceiveNewEvent
	// and the handler for getting a session (only one from the user's perspective)
	mux[RestResource{Route: "/session", Method: http.MethodGet}] = handlers.GetSession
}

// starting up the server
func setupServer() {

	port := ":8000"
	// equivalent of http.ListenAndServer but not using the builtin server
	server := http.Server{
		Addr:    port,
		Handler: &customHandler{},
	}

	log.Println("Server running on port", port)
	log.Fatal(server.ListenAndServe())
}

func main() {

	log.Println("Server starting ...")
	setupHandlers()
	setupServer()

}
