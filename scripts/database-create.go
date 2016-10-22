package main

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

var (
	pool *pgx.ConnPool
)

func main() {
	fmt.Println("Create database" + os.Getenv("DATABASE_URL"))

	// Connect to Postgres database
	conn, err := pgx.ParseURI(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Invalid DATABASE_URL format", err)
		os.Exit(1)
	}

	// Create connection pool
	pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: 2,
	})
	if err != nil {
		fmt.Println("Unable to create connection pool", err)
		os.Exit(1)
	}

	createAccount()
	createToken()
}

func createToken() {
	// Search account tokens for parameter
	rows, err := pool.Query(`CREATE TABLE token (
    created timestamp DEFAULT current_timestamp NOT NULL,
    token character varying(100) NOT NULL,
    account integer NOT NULL
  );`)
	defer rows.Close()

	if err != nil {
		fmt.Println("Error select", err)
		return
	}
}

func createAccount() {
	// Search account tokens for parameter
	rows, err := pool.Query(`CREATE TABLE account (
    id serial primary key,
    address varchar(100) NOT NULL,
    verified boolean DEFAULT false NOT NULL,
    created timestamp DEFAULT current_timestamp NOT NULL,
    token character varying(100)
  );`)
	defer rows.Close()

	if err != nil {
		fmt.Println("Error select", err)
		return
	}
}
