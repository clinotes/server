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
		var data APIRequestStructCreateToken
		if ensureJSONPayload(req, res, &data) != nil {
			return
		}

		// Generate new random token
		token := random(36)
		hashed, _ := passlib.Hash(token)

		// Lookup account ID
		accountID, err := accountIDByAddress(data.Address)
		if err != nil {
			writeJSONError(res, "Unknown account address")
			return
		}

		// Save token to database
		if _, err := pool.Exec("addToken", hashed, accountID); err == nil {
			// Send confirmation mail using Postmark
			pmark.SendTemplatedEmail(postmark.TemplatedEmail{
				TemplateId: int64(postmarkTemplateIDauthTokenCreate),
				TemplateModel: map[string]interface{}{
					"token": token,
				},
				From:    "mail@clinot.es",
				To:      data.Address,
				ReplyTo: "\"CLINotes\" <mail@clinot.es>",
			})

			writeJSONResponse(res)
			return
		}

		writeJSONError(res, "Failed to create token for account")
	},
}
