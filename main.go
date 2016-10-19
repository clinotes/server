package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/keighl/postmark"
)

var pool *pgx.ConnPool
var pmark *postmark.Client

func init() {
	// Create Postmark API client
	pmark = postmark.NewClient(os.Getenv("POSTMARK_API_KEY"), "")
	// Connect to Postgres database
	conn, err := pgx.ParseURI(os.Getenv("DATABASE_URL"))
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: 5,
		AfterConnect:   registerQueries,
	}

	// Create connection pool
	pool, err = pgx.NewConnPool(connPoolConfig)
	if err != nil {
		fmt.Println("Unable to create connection pool", err)
		os.Exit(1)
	}
}

func main() {
	router := mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()

	api.HandleFunc(APIRouteCreateToken.URL, APIRouteCreateToken.Handler)
	api.HandleFunc(APIRouteCreateUser.URL, APIRouteCreateUser.Handler)
	api.HandleFunc(APIRouterAuth.URL, APIRouterAuth.Handler)
	api.HandleFunc(APIRouterVerifyUser.URL, APIRouterVerifyUser.Handler)

	// Listen on PORT only on non-local environment
	if os.Getenv("ENV") == "local" {
		http.ListenAndServe("localhost:8000", router)
	} else {
		http.ListenAndServe(":"+os.Getenv("PORT"), router)
	}
}
