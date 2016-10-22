package setup

import (
	"fmt"

	"github.com/jackc/pgx"
)

// Run creates the database structure if needed
func Run(pool *pgx.Conn) {
	createAccount(pool)
	createToken(pool)
	createNote(pool)
}

func createNote(pool *pgx.Conn) {
	// Search account tokens for parameter
	rows, err := pool.Query(`CREATE TABLE note (
		account integer NOT NULL,
		text character varying(250) unique NOT NULL,
		created timestamp DEFAULT current_timestamp NOT NULL
  );`)
	defer rows.Close()

	if err != nil {
		fmt.Println("Error select", err)
		return
	}
}

func createToken(pool *pgx.Conn) {
	// Search account tokens for parameter
	rows, err := pool.Query(`CREATE TABLE token (
		account integer NOT NULL,
		token character varying(100) unique NOT NULL,
		created timestamp DEFAULT current_timestamp NOT NULL
  );`)
	defer rows.Close()

	if err != nil {
		fmt.Println("Error select", err)
		return
	}
}

func createAccount(pool *pgx.Conn) {
	// Search account tokens for parameter
	rows, err := pool.Query(`CREATE TABLE account (
    id serial primary key,
    address varchar(100) unique NOT NULL,
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
