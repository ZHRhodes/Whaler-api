package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/heroku/whaler-api/app"
	"github.com/heroku/whaler-api/controllers"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)

	router.HandleFunc("/api/user/create", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	//this endpoint should probably eventually be access controlled if not removed
	router.HandleFunc("/api/org/create", controllers.CreateOrg).Methods("POST")

	router.HandleFunc("/api/org", controllers.FetchOrg).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
