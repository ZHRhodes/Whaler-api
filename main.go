package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/heroku/whaler-api/graph"
	"github.com/heroku/whaler-api/graph/generated"
)

func main() {
	// router := mux.NewRouter()
	// router.Use(app.JwtAuthentication)

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

	//GraphQL
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/schema", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	//End GraphQL

	fmt.Println(port)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Print(err)
	}
}
