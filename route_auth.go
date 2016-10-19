package main

import (
	"fmt"
	"net/http"

	"gopkg.in/hlandau/passlib.v1"
)

// APIRequestStructAuth is
type APIRequestStructAuth struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// APIRouterAuth is
var APIRouterAuth = Route{
	"/auth",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var data APIRequestStructAuth
		if ensureJSONPayload(req, res, &data) != nil {
			return
		}

		// Lookup account ID
		accountID, err := accountIDByAddress(data.Address)
		if err != nil {
			writeJSONError(res, "Unknown account address")
			return
		}

		// Search account tokens for parameter
		rows, _ := pool.Query("select token from token WHERE account = $1", accountID)
		for rows.Next() {
			values, err := rows.Values()
			// Error with result
			if err != nil {
				writeJSONResponse(res)
				return
			}

			_, err = passlib.Verify(data.Token, fmt.Sprintf("%s", values[0]))
			// Found matching token
			if err == nil {
				writeJSONResponse(res)
				return
			}
		}

		// Return error if no token is found
		writeJSONError(res, "Unknown account token")
	},
}
