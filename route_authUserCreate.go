package main

import (
	"net/http"

	"github.com/keighl/postmark"
	"gopkg.in/hlandau/passlib.v1"
)

// APIRequestStructCreateUser is
type APIRequestStructCreateUser struct {
	Address string `json:"address"`
}

var postmarkTemplateIDauthUserCreate = 1012641

// APIRouteCreateUser is
var APIRouteCreateUser = Route{
	"/auth/user/create",
	func(res http.ResponseWriter, req *http.Request) {
		var data APIRequestStructCreateUser
		if ensureJSONPayload(req, res, &data) != nil {
			return
		}

		// Create token for user verification
		verify := random(36)
		hashed, _ := passlib.Hash(verify)

		// Save account to database
		if _, err := pool.Exec("addUser", data.Address, hashed); err == nil {
			// Send confirmation mail using Postmark
			pmark.SendTemplatedEmail(postmark.TemplatedEmail{
				TemplateId: int64(postmarkTemplateIDauthUserCreate),
				TemplateModel: map[string]interface{}{
					"token": verify,
				},
				From:    "mail@clinot.es",
				To:      data.Address,
				ReplyTo: "\"CLINotes\" <mail@clinot.es>",
			})

			writeJSONResponse(res)
			return
		}

		writeJSONError(res, "Unable to create account")
	},
}
