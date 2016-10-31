package main

import (
	"fmt"

	"github.com/clinotes/server/data"
	"github.com/clinotes/server/setup"
	"github.com/jackc/pgx"
)

func prepareQueries(conn *pgx.Conn, list map[string]string) error {
	for name, query := range list {
		fmt.Println("Register query:", name)
		_, err := conn.Prepare(name, query)

		if err != nil {
			fmt.Println("Failed to prepare query.", err)
			return err
		}
	}

	return nil
}

func registerQueries(conn *pgx.Conn) error {
	setup.Run(conn)
	fmt.Println("Ran setup")

	queryList := []map[string]string{
		data.AccountQueries,
		data.TokenQueries,
		data.SubscriptionQueries,
	}

	for _, item := range queryList {
		err := prepareQueries(conn, item)
		if err != nil {
			return err
		}
	}

	return nil
}
