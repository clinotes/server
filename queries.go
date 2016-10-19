package main

import (
	"fmt"

	"github.com/jackc/pgx"
)

var queries = map[string]string{
	"addUser": `
		insert into account (address, token, verified)
		values($1, $2, false)
	`,
	"addToken": `
    insert into token (token, account)
    values($1, $2)
	`,
	"getUser": `
    select id from account where address=$1 AND verified = TRUE
	`,
	"getUnverifiedUser": `
    select token from account where address=$1 AND verified = FALSE
	`,
	"verifyUser": `
    update account
    set token = NULL, verified = TRUE
    WHERE address = $1
	`,
}

func registerQueries(conn *pgx.Conn) error {
	for name, query := range queries {
		_, err := conn.Prepare(name, query)

		if err != nil {
			fmt.Println("Failed to prepare query.", err)
			return err
		}
	}

	return nil
}
