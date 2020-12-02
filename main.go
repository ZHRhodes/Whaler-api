package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"
	"github.com/heroku/whaler-api/controllers"
	"github.com/heroku/whaler-api/graph"
	"github.com/heroku/whaler-api/graph/generated"
	"github.com/heroku/whaler-api/middleware"
	"github.com/heroku/whaler-api/models"
)

func main() {
	router := mux.NewRouter()
	router.Use(middleware.JwtAuthentication)
	// router.Use(middleware.ParseUserIDFromToken)

	resolver := graph.Resolver{DB: models.DB()}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}))
	router.Handle("/query", srv)
	router.Handle("/schema", playground.Handler("GraphQL playground", "/query"))

	router.HandleFunc("/api/user/login", controllers.Authenticate)
	router.HandleFunc("/api/user/refresh", controllers.Refresh)

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

	router.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			fmt.Print(fmt.Sprint("Websocket error 1 ", err))
		}
		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					fmt.Print(fmt.Sprint("Websocket error 2 ", err))
				}
				fmt.Print(fmt.Sprint("The message is ", msg))
				fmt.Print(fmt.Sprint("The op is", op))
				err = wsutil.WriteServerMessage(conn, op, msg)
				if err != nil {
					fmt.Print(fmt.Sprint("Websocket error 3 ", err))
				}
			}
		}()
	})

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
