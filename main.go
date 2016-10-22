package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/clinotes/server/route"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/keighl/postmark"
)

var (
	pool  *pgx.ConnPool
	pmark *postmark.Client

	maxDBConnections int

	httpHostname string
	httpPort     string

	router *mux.Router
)

func init() {
	// Read value for max DB connections from ENV
	maxDBConnections, err := strconv.Atoi(os.Getenv("MAX_DB_CONNECTIONS"))
	if err != nil {
		fmt.Println("Failed to set maxDBConnections", err)
		os.Exit(1)
	}

	// Connect to Postgres database
	conn, err := pgx.ParseURI(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Invalid DATABASE_URL format", err)
		os.Exit(1)
	}

	// Create connection pool
	pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: maxDBConnections,
		AfterConnect:   registerQueries,
	})
	if err != nil {
		fmt.Println("Unable to create connection pool", err)
		os.Exit(1)
	}

	// Create Postmark API client
	pmark = postmark.NewClient(os.Getenv("POSTMARK_API_KEY"), "")

	// Create mux router
	router = mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()

	// Configure path handlers
	for _, r := range route.List(pool, pmark) {
		api.HandleFunc(r.URL, r.Handler).Methods("POST")
	}
}

func main() {
	// Check if running on local environment and set hostname to avoid
	// annoying MacOS security warnings.
	if os.Getenv("ENV") == "local" {
		httpHostname = "localhost"
		httpPort = "8000"
	} else {
		httpHostname = ""
		httpPort = os.Getenv("PORT")
	}

	// Listen on PORT only on non-local environment
	fmt.Printf("Started CLInotes API endpoint on port %s\n", httpPort)
	http.ListenAndServe(httpHostname+":"+httpPort, router)
}
