package data

import (
	"flag"
	"fmt"
	"os"
	"testing"

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

func TestMain(m *testing.M) {
	// conn, _ := pgx.ParseURI(os.Getenv("DATABASE_URL"))
	pool, _ = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgx.ConnConfig{Host: "127.0.0.1", User: "sebastian", Password: "", Database: "testing"},
		MaxConnections: 1,
		AfterConnect: func(conn *pgx.Conn) error {
			queryList := []map[string]string{
				AccountQueries,
				TokenQueries,
				SubscriptionQueries,
			}

			for _, item := range queryList {
				err := prepareQueries(conn, item)
				if err != nil {
					return err
				}
			}

			return nil
		},
	})

	Pool(pool)

	flag.Parse()
	os.Exit(m.Run())
}