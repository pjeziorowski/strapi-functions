package main

import (
	"functions/backend/config"
	"functions/backend/handler/auth/signin"
	"functions/backend/handler/auth/signup"
	"functions/backend/handler/create_project"
	"functions/backend/handler/delete_project"
	"functions/backend/handler/start_project"
	"functions/backend/handler/stop_project"
	"log"
	"net/http"
)

// HTTP server for the handler
func main() {
	// check server configuration
	configErrors := config.CheckServerConfig()
	if configErrors != nil {
		for _, err := range configErrors {
			log.Printf(err.Error())
		}
		log.Fatal("killing the server")
	}

	mux := http.NewServeMux()

	// register handlers
	mux.HandleFunc("/signin", signin.Handler)
	mux.HandleFunc("/signup", signup.Handler)
	mux.HandleFunc("/create_project", create_project.Handler)
	mux.HandleFunc("/start_project", start_project.Handler)
	mux.HandleFunc("/stop_project", stop_project.Handler)
	mux.HandleFunc("/delete_project", delete_project.Handler)

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
