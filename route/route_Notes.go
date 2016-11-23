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

// APIRequestStructNotes is
type APIRequestStructNotes struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// APIResponseStructNote is
type APIResponseStructNote struct {
	Text    string
	Created time.Time
}

// APIRouteNotes is
var APIRouteNotes = Route{
	"/notes",
	func(res http.ResponseWriter, req *http.Request) (error, interface{}) {
		// Parse JSON request
		var reqData APIRequestStructNotes
		if err := checkJSONBody(req, res, &reqData); err != nil {
			return err, nil
		}

		// Get account
		account, err := data.AccountByAddress(reqData.Address)
		if err != nil {
			return errors.New("Unknown account address"), nil
		}

		if !account.IsVerified() {
			return errors.New("Account not verified"), nil
		}

		// Check if account has requested token
		_, err = account.GetToken(reqData.Token, data.TokenTypeAccess)
		if err != nil {
			return errors.New("Unable to use provided token"), nil
		}

		list, err := data.NoteListByAccount(account.ID())
		if err != nil {
			return errors.New("Failed to get notes"), nil
		}

		var noteList []APIResponseStructNote
		for i := 0; i < len(list); i++ {
			noteList = append(noteList, APIResponseStructNote{list[i].Text(), list[i].CreatedOn()})
		}

		return nil, noteList
	},
}
