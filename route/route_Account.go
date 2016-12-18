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
	"errors"
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
	func(res http.ResponseWriter, req *http.Request) (interface{}, error) {
		// Parse JSON request
		var reqData APIRequestStructMe
		if err := checkJSONBody(req, res, &reqData); err != nil {
			return nil, err
		}

		// Get account
		account, err := data.AccountByAddress(reqData.Address)
		if err != nil {
			return nil, errors.New("Unknown account address")
		}

		// Check if account has requested token
		_, err = account.GetToken(reqData.Token, data.TokenTypeAccess)
		if err != nil {
			return nil, errors.New("Unable to use provided token")
		}

		// Verify account
		account, err = account.Verify()
		if err != nil {
			return nil, errors.New("Unable to use provided token")
		}

		return APIResponseStructAccount{
			account.Address,
			account.Created,
			account.HasSubscription(),
		}, nil
	},
}
