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
