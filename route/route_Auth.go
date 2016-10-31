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

	"github.com/clinotes/server/data"
)

// APIRequestStructAuth is
type APIRequestStructAuth struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// APIRouteAuth is
var APIRouteAuth = Route{
	"/auth",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var reqData APIRequestStructAuth
		if ensureJSONPayload(req, res, &reqData) != nil {
			return
		}

		// Get account
		account, err := data.AccountByAddress(reqData.Address)
		if err != nil {
			writeJSONError(res, "Unknown account address")
			return
		}

		if !account.IsVerified() {
			writeJSONError(res, "Account not verified")
			return
		}

		// Check if account has requested token
		_, err = account.GetToken(reqData.Token, data.TokenTypeAccess)
		if err != nil {
			writeJSONError(res, "Unable to use provided token")
			return
		}

		writeJSONResponse(res)
	},
}