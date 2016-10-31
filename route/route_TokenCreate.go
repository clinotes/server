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
