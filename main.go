package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/heroku/whaler-api/controllers"
	"github.com/heroku/whaler-api/graph"
	"github.com/heroku/whaler-api/graph/generated"
	"github.com/heroku/whaler-api/middleware"
)

func main() {
	router := mux.NewRouter()
	router.Use(middleware.ParseUserIDFromToken)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/query", srv)
	router.Handle("/schema", playground.Handler("GraphQL playground", "/query"))

	http.HandleFunc("/api/user/login", controllers.Authenticate)
	http.HandleFunc("/api/user/refresh", controllers.Refresh)

	// router.HandleFunc("/api/user/create", controllers.CreateUser).Methods("POST")
	// router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	// router.HandleFunc("/api/user/refresh", controllers.Refresh).Methods("POST")
	// router.HandleFunc("/api/user/logout", controllers.LogOut).Methods("POST")

	// router.HandleFunc("/api/org/create", controllers.CreateOrg).Methods("POST")
	// router.HandleFunc("/api/workspace/create", controllers.CreateWorkspace).Methods("POST")
	// router.HandleFunc("/api/account/create", controllers.CreateAccount).Methods("POST")
	// router.HandleFunc("/api/contact/create", controllers.CreateContact).Methods("POST")

	// router.HandleFunc("/api/org", controllers.FetchOrg).Methods("GET")
	// router.HandleFunc("/api/workspace", controllers.FetchWorkspace).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("connect to port %s for GraphQL playground", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
