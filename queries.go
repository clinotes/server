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

package main

import (
	"fmt"

	"github.com/clinotes/server/data"
	"github.com/clinotes/server/setup"
	"github.com/jackc/pgx"
)

func prepareQueries(conn *pgx.Conn, list map[string]string) error {
	for name, query := range list {
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
