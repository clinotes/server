package main

import (
	"net/http"

	"github.com/keighl/postmark"
	"gopkg.in/hlandau/passlib.v1"
)

// APIRequestStructCreateToken is
type APIRequestStructCreateToken struct {
	Address string `json:"address"`
}

var postmarkTemplateIDauthTokenCreate = 1010802

// APIRouteCreateToken is
var APIRouteCreateToken = Route{
	"/auth/token/create",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var data APIRequestStructCreateToken
		if ensureJSONPayload(req, res, &data) != nil {
			return
		}

		// Generate new random token
		token := random(36)
		hashed, err := passlib.Hash(token)

		if err != nil {
			writeJSONError(res, "Unable to create token")
			return
		}

		// Lookup account ID
		accountID, err := accountIDByAddress(data.Address)
		if err != nil {
			writeJSONError(res, "Unknown account address")
			return
		}

		// Save token to database
		if _, err := pool.Exec("addToken", hashed, accountID); err == nil {
			// Send confirmation mail using Postmark
			_, err := pmark.SendTemplatedEmail(postmark.TemplatedEmail{
				TemplateId: int64(postmarkTemplateIDauthTokenCreate),
				TemplateModel: map[string]interface{}{
					"token": token,
				},
				From:    "mail@clinot.es",
				To:      data.Address,
				ReplyTo: "\"CLINotes\" <mail@clinot.es>",
			})

			if err != nil {
				writeJSONError(res, "Unable to create token for account")
				return
			}

			writeJSONResponse(res)
			return
		}

		writeJSONError(res, "Failed to create token for account")
	},
}
