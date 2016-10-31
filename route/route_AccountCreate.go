package route

import (
	"net/http"

	"github.com/clinotes/server/data"
	"github.com/keighl/postmark"
)

// APIRequestStructCreateUser is
type APIRequestStructCreateUser struct {
	Address string `json:"address"`
}

var postmarkTemplateIDauthUserCreate = 1012641

// APIRouteAccountCreate is
var APIRouteAccountCreate = Route{
	"/account/create",
	func(res http.ResponseWriter, req *http.Request) {
		var reqData APIRequestStructCreateUser
		if ensureJSONPayload(req, res, &reqData) != nil {
			return
		}

		account := data.AccountNew(reqData.Address)
		account, err := account.Store()

		// If account cannot be created, fail
		if err != nil {
			writeJSONError(res, "Unable to create account")
			return
		}

		token := data.TokenNew(account.ID(), data.TokenTypeMaintenace)
		tokenRaw := token.Raw()
		token, err = token.Store()

		// If token cannot be created, fail and remove user
		if err != nil {
			account.Remove()
			writeJSONError(res, "Unable to create account")
			return
		}

		// Send confirmation mail using Postmark
		_, err = pmark.SendTemplatedEmail(postmark.TemplatedEmail{
			TemplateId: int64(postmarkTemplateIDauthUserCreate),
			TemplateModel: map[string]interface{}{
				"token": tokenRaw,
			},
			From:    "mail@clinot.es",
			To:      account.Address(),
			ReplyTo: "\"CLINotes\" <mail@clinot.es>",
		})

		// If mail cannot be sent, fail and remove user
		if err != nil {
			account.Remove()
			writeJSONError(res, "Unable to create account")
			return
		}

		// Done!
		writeJSONResponse(res)
	},
}
