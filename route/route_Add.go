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

	"github.com/clinotes/server/data"
)

// APIRequestStructAdd is
type APIRequestStructAdd struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Note    string `json:"note"`
}

// APIRouteAdd is
var APIRouteAdd = Route{
	"/add",
	func(res http.ResponseWriter, req *http.Request) (interface{}, error) {
		// Parse JSON request
		var reqData APIRequestStructAdd
		if err := checkJSONBody(req, res, &reqData); err != nil {
			return nil, err
		}

		// Get account
		account, err := data.AccountByAddress(reqData.Address)
		if err != nil {
			return nil, errors.New("Unknown account address")
		}

		if !account.IsVerified() {
			return nil, errors.New("Account not verified")
		}

		// Check if account has requested token
		_, err = account.GetToken(reqData.Token, data.TokenTypeAccess)
		if err != nil {
			return nil, errors.New("Unable to use provided token")
		}

		note := data.NoteNew(account.ID(), reqData.Note)
		note, err = note.Store()

		if err != nil {
			return nil, errors.New("Unable to store note")
		}

		return nil, nil
	},
}
