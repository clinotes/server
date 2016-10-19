package main

import (
	"net/http"

	"github.com/keighl/postmark"
	"gopkg.in/hlandau/passlib.v1"
)

// APIRequestStructVerifyUser is
type APIRequestStructVerifyUser struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

var postmarkTemplateIDauthUserVerify = 1012661

// APIRouterVerifyUser is
var APIRouterVerifyUser = Route{
	"/auth/user/verify",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var data APIRequestStructVerifyUser
		if ensureJSONPayload(req, res, &data) != nil {
			return
		}

		// Lookup account ID
		hash, err := tokenByUnverifiedAddress(data.Address)
		if err != nil {
			writeJSONError(res, "Unknown account address")
			return
		}

		// Verify token
		_, err = passlib.Verify(data.Token, hash)
		if err != nil {
			writeJSONError(res, "Invalid verification token")
			return
		}

		if _, err = pool.Exec("verifyUser", data.Address); err == nil {
			// Send confirmation mail using Postmark
			_, err := pmark.SendTemplatedEmail(postmark.TemplatedEmail{
				TemplateId: int64(postmarkTemplateIDauthUserVerify),
				From:       "mail@clinot.es",
				To:         data.Address,
				TemplateModel: map[string]interface{}{
					"token": data.Token,
				},
				ReplyTo: "\"CLINotes\" <mail@clinot.es>",
			})

			if err != nil {
				writeJSONError(res, "Unable to verify account")
				return
			}

			writeJSONResponse(res)
			return
		}

		writeJSONError(res, "Failed to verify account")
	},
}
