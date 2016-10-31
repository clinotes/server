package route

import (
	"net/http"

	"github.com/clinotes/server/data"
)

// APIRequestStructAuth is
type APIRequestStructAuth struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// APIRouteAuth is
var APIRouteAuth = Route{
	"/auth",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var reqData APIRequestStructAuth
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

		writeJSONResponse(res)
	},
}
