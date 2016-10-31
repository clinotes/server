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

// APIRequestStructCreateToken is
type APIRequestStructCreateToken struct {
	Address string `json:"address"`
}

var postmarkTemplateIDauthTokenCreate = 1010802

// APIRouteTokenCreate is
var APIRouteTokenCreate = Route{
	"/token/create",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var reqData APIRequestStructCreateToken
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

		token := data.TokenNew(account.ID(), data.TokenTypeAccess)
		tokenRaw := token.Raw()
		token, err = token.Store()

		_, err = pmark.SendTemplatedEmail(postmark.TemplatedEmail{
			TemplateId: int64(postmarkTemplateIDauthTokenCreate),
			TemplateModel: map[string]interface{}{
				"token": tokenRaw,
			},
			From:    "mail@clinot.es",
			To:      reqData.Address,
			ReplyTo: "\"CLINotes\" <mail@clinot.es>",
		})

		if err != nil {
			token.Remove()
			writeJSONError(res, "Unable to create token for account")
			return
		}

		writeJSONResponse(res)
	},
}
