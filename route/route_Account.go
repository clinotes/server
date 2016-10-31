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

package route

import (
	"net/http"
	"time"

	"github.com/clinotes/server/data"
)

// APIRequestStructMe is
type APIRequestStructMe struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// APIResponseStructAccount is
type APIResponseStructAccount struct {
	Address      string
	Created      time.Time
	Subscription bool
}

// APIRouteAccount is
var APIRouteAccount = Route{
	"/account",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var reqData APIRequestStructMe
		if ensureJSONPayload(req, res, &reqData) != nil {
			return
		}

		// Get account
		account, err := data.AccountByAddress(reqData.Address)
		if err != nil {
			writeJSONError(res, "Unknown account address")
			return
		}

		// Check if account has requested token
		_, err = account.GetToken(reqData.Token, data.TokenTypeAccess)
		if err != nil {
			writeJSONError(res, "Unable to use provided token")
			return
		}

		// Verify account
		account, err = account.Verify()
		if err != nil {
			writeJSONError(res, "Unable to use provided token")
			return
		}

		resData := APIResponseStructAccount{
			account.Address(),
			account.CreatedOn(),
			account.HasSubscription(),
		}

		writeJSONResponseData(res, resData)
	},
}
