package route

import (
	"net/http"
	"time"

	"github.com/clinotes/server/data"
)

// APIRequestStructMe is
type APIRequestStructMe struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// APIResponseStructAccount is
type APIResponseStructAccount struct {
	Address      string
	Created      time.Time
	Subscription bool
}

// APIRouteAccount is
var APIRouteAccount = Route{
	"/account",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var reqData APIRequestStructMe
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
		_, err = account.GetToken(reqData.Token, data.TokenTypeAccess)
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

		resData := APIResponseStructAccount{
			account.Address(),
			account.CreatedOn(),
			account.HasSubscription(),
		}

		writeJSONResponseData(res, resData)
	},
}
