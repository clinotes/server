/**
 * clinot.es server
 * Copyright (C) 2016 Sebastian Müller
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.

 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package data

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/jmoiron/sqlx"
)

func prepareQueries(conn *pgx.Conn, list map[string]string) error {
	Setup(conn)

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
	db, _ = sqlx.Open("pgx", os.Getenv("DATABASE_URL"))
	Database(db)

	flag.Parse()
	os.Exit(m.Run())
}
