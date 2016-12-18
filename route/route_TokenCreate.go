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

// APIRequestStructCreateToken is
type APIRequestStructCreateToken struct {
	Address string `json:"address"`
}

// APIRouteTokenCreate is
var APIRouteTokenCreate = Route{
	"/token/create",
	func(res http.ResponseWriter, req *http.Request) (interface{}, error) {
		// Parse JSON request
		var reqData APIRequestStructCreateToken
		if err := checkJSONBody(req, res, &reqData); err != nil {
			return nil, err
		}

		// Get account
		account, err := data.AccountByAddress(reqData.Address)
		if err != nil {
			return nil, errors.New("Unknown account address")
		}

		if !account.Verified {
			return nil, errors.New("Account not verified")
		}

		token := data.TokenNew(account.ID, data.TokenTypeAccess)
		tokenRaw := token.Raw()
		token, err = token.Store()

		_, err = sendTokenWithTemplate(reqData.Address, tokenRaw, conf.TemplateToken)
		if err != nil {
			token.Remove()
			return nil, errors.New("Unable to create token for account")
		}

		return nil, nil
	},
}
