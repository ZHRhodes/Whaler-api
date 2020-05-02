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
	router.HandleFunc("/api/user/refresh", controllers.Refresh).Methods("POST")
	router.HandleFunc("/api/user/logout", controllers.LogOut).Methods("POST")

	router.HandleFunc("/api/org/create", controllers.CreateOrg).Methods("POST")
	router.HandleFunc("/api/workspace/create", controllers.CreateWorkspace).Methods("POST")
	router.HandleFunc("/api/account/create", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/contact/create", controllers.CreateContact).Methods("POST")

	router.HandleFunc("/api/org", controllers.FetchOrg).Methods("GET")
	router.HandleFunc("/api/workspace", controllers.FetchWorkspace).Methods("GET")

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
