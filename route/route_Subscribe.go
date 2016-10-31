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
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/clinotes/server/data"
	stripe "github.com/stripe/stripe-go"
	stripeCustomer "github.com/stripe/stripe-go/customer"
	stripeSub "github.com/stripe/stripe-go/sub"
	stripeToken "github.com/stripe/stripe-go/token"
)

// APIRequestStructSubscribe is
type APIRequestStructSubscribe struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Number  string `json:"number"`
	Expire  string `json:"expire"`
	CVC     string `json:"cvc"`
}

// APIRouteSubscribe is
var APIRouteSubscribe = Route{
	"/subscribe",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var reqData APIRequestStructSubscribe
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

		stripe.Key = os.Getenv("STRIPE_API_KEY")
		expire := strings.Split(reqData.Expire, "/")
		t, err := stripeToken.New(&stripe.TokenParams{
			Card: &stripe.CardParams{
				Number: reqData.Number,
				Month:  expire[0],
				Year:   expire[1],
				CVC:    reqData.CVC,
			},
		})

		if err != nil {
			writeJSONError(res, "Invalid card information")
			return
		}

		customerParams := &stripe.CustomerParams{
			Desc: fmt.Sprintf("%s (#%d)", account.Address(), account.ID()),
		}
		customerParams.SetSource(t.ID)
		c, err := stripeCustomer.New(customerParams)

		if err != nil {
			writeJSONError(res, "Invalid account information")
			return
		}

		s, err := stripeSub.New(&stripe.SubParams{
			Customer: c.ID,
			Plan:     "blue",
		})

		if err != nil {
			writeJSONError(res, "Invalid account information")
			return
		}

		subscription := data.SubscriptionNew(account.ID(), s.ID)
		subscription, err = subscription.Store()
		subscription, err = subscription.Activate()

		if err != nil {
			writeJSONError(res, "Invalid account information")
			return
		}

		writeJSONResponse(res)
	},
}
