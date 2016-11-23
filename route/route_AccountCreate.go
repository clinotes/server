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

// APIRequestStructCreateUser is
type APIRequestStructCreateUser struct {
	Address string `json:"address"`
}

// APIRouteAccountCreate is
var APIRouteAccountCreate = Route{
	"/account/create",
	func(res http.ResponseWriter, req *http.Request) (error, interface{}) {
		var reqData APIRequestStructCreateUser
		if err := checkJSONBody(req, res, &reqData); err != nil {
			return err, nil
		}

		account := data.AccountNew(reqData.Address)
		account, err := account.Store()

		// If account cannot be created, fail
		if err != nil {
			return errors.New("Unable to create account"), nil
		}

		token := data.TokenNew(account.ID(), data.TokenTypeMaintenace)
		tokenRaw := token.Raw()
		token, err = token.Store()

		// If token cannot be created, fail and remove user
		if err != nil {
			account.Remove()
			return errors.New("Unable to create account"), nil
		}

		// Send confirmation mail using Postmark
		_, err = sendTokenWithTemplate(account.Address(), tokenRaw, conf.TemplateWelcome)
		if err != nil {
			return errors.New("Unable to send welcome mail"), nil
		}

		// If mail cannot be sent, fail and remove user
		if err != nil {
			account.Remove()
			return errors.New("Unable to create account"), nil
		}

		// Done!
		return nil, nil
	},
}
