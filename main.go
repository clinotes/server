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
	"net/http"
	"os"

	"github.com/clinotes/server/data"
	"github.com/clinotes/server/route"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/keighl/postmark"
	"github.com/spf13/viper"
)

var (
	version                = "0.0.6"
	versionClientSupported = "0.2.0"

	pool  *pgx.ConnPool
	pmark *postmark.Client

	maxDBConnections int
	connectionURL    string
	postmarkAPIToken string

	httpHostname string
	httpPort     string

	postmarkTemplateIDWelcome int64
	postmarkTemplateIDConfirm int64
	postmarkTemplateIDToken   int64
	postmarkFrom              string
	postmarkReplyTo           string

	router *mux.Router
)

func readEnvironment() {
	viper.AutomaticEnv()
	viper.ReadInConfig()

	maxDBConnections = viper.GetInt("MAX_DB_CONNECTIONS")
	connectionURL = viper.GetString("DATABASE_URL")
	postmarkAPIToken = viper.GetString("POSTMARK_API_KEY")

	postmarkTemplateIDWelcome = viper.GetInt64("POSTMARK_TEMPLATE_WELCOME")
	postmarkTemplateIDConfirm = viper.GetInt64("POSTMARK_TEMPLATE_CONFIRM")
	postmarkTemplateIDToken = viper.GetInt64("POSTMARK_TEMPLATE_TOKEN")
	postmarkFrom = viper.GetString("POSTMARK_FROM")
	postmarkReplyTo = viper.GetString("POSTMARK_REPLY_TO")
}

func checkEnvironment() {
	if maxDBConnections <= 0 {
		fmt.Println("Please set MAX_DB_CONNECTIONS > 0")
		os.Exit(1)
	}

	if postmarkTemplateIDWelcome <= 0 {
		fmt.Println("Please set POSTMARK_TEMPLATE_WELCOME > 0")
		os.Exit(1)
	}

	if postmarkTemplateIDConfirm <= 0 {
		fmt.Println("Please set POSTMARK_TEMPLATE_CONFIRM > 0")
		os.Exit(1)
	}

	if postmarkTemplateIDToken <= 0 {
		fmt.Println("Please set POSTMARK_TEMPLATE_TOKEN > 0")
		os.Exit(1)
	}

	if postmarkAPIToken == "" {
		fmt.Println("Please set POSTMARK_API_KEY")
		os.Exit(1)
	}

	if postmarkFrom == "" {
		fmt.Println("Please set POSTMARK_FROM")
		os.Exit(1)
	}

	if postmarkReplyTo == "" {
		fmt.Println("Please set POSTMARK_REPLY_TO")
		os.Exit(1)
	}
}

func connect() *pgx.ConnPool {
	conn, err := pgx.ParseURI(connectionURL)
	if err != nil {
		fmt.Println("Invalid DATABASE_URL format", err)
		os.Exit(1)
	}

	// Create connection pool
	pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: maxDBConnections,
		AfterConnect: func(conn *pgx.Conn) error {
			data.Setup(conn)

			for _, item := range data.Queries {
				if queryError := prepareQueries(conn, item); queryError != nil {
					return queryError
				}
			}

			return nil
		},
	})

	if err != nil {
		fmt.Println("Unable to create connection pool", err)
		os.Exit(1)
	}

	return pool
}

func init() {
	readEnvironment()
	checkEnvironment()

	// Connect to data pool
	data.Pool(connect())

	// Create mux router
	router = mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()

	config := route.Configuration{
		TemplateWelcome: postmarkTemplateIDWelcome,
		TemplateConfirm: postmarkTemplateIDConfirm,
		TemplateToken:   postmarkTemplateIDToken,
		PostmarkToken:   postmarkAPIToken,
		PostmarkFrom:    postmarkFrom,
		PostmarkReplyTo: postmarkReplyTo,
	}

	// Configure path handlers
	for _, r := range route.Routes(config) {
		api.Handle(r.URL, route.Handler(r.Handler)).Methods("POST")
	}

	api.HandleFunc(
		"/version",
		func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json; charset=utf-8")
			res.Write([]byte(`{"version":"` + version + `","client":"` + versionClientSupported + `"}`))
		},
	)
}

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
