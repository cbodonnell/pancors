package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/michaljanocko/pancors"
	"github.com/michaljanocko/pancors/pkg/middleware"
)

func getAllowOrigin() string {
	if origin, ok := os.LookupEnv("ALLOW_ORIGIN"); ok {
		return origin
	}
	return "*"
}

func getAllowCredentials() string {
	if credentials, ok := os.LookupEnv("ALLOW_CREDENTIALS"); ok {
		return credentials
	}
	return "true"
}

func getListenPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return port
	}
	return "8080"
}

func getAuthEndpoint() string {
	if endpoint, ok := os.LookupEnv("AUTH_ENDPOINT"); ok {
		return endpoint
	}
	return ""
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", pancors.HandleProxyWith(getAllowOrigin(), getAllowCredentials()))

	authEndpoint := getAuthEndpoint()
	if authEndpoint != "" {
		log.Printf("Authenticating with %s", authEndpoint)
		auth := middleware.NewAuthMiddleware(authEndpoint)
		r.Use(auth)
	}

	// TODO: CORS in development to allow unauthorized testing

	port := getListenPort()
	log.Printf("PanCORS started listening on %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
