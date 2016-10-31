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
	"github.com/keighl/postmark"
)

// APIRequestStructVerifyUser is
type APIRequestStructVerifyUser struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

var postmarkTemplateIDauthUserVerify = 1012661

// APIRouteAccountVerify is
var APIRouteAccountVerify = Route{
	"/account/verify",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var reqData APIRequestStructVerifyUser
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
		_, err = account.GetToken(reqData.Token, data.TokenTypeMaintenace)
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

		// Send confirmation mail
		pmark.SendTemplatedEmail(postmark.TemplatedEmail{
			TemplateId: int64(postmarkTemplateIDauthUserVerify),
			From:       "mail@clinot.es",
			To:         account.Address(),
			TemplateModel: map[string]interface{}{
				"token": reqData.Token,
			},
			ReplyTo: "\"CLINotes\" <mail@clinot.es>",
		})

		writeJSONResponse(res)
	},
}
