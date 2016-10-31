package route

import (
	"fmt"
	"net/http"

	"github.com/clinotes/server/data"
)

// APIRequestStructAdd is
type APIRequestStructAdd struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Note    string `json:"note"`
}

// APIRouteAdd is
var APIRouteAdd = Route{
	"/add",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var reqData APIRequestStructAdd
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

		fmt.Printf("Add Note: %+v\n", data.NoteNew(reqData.Note))

		// Return error if no token is found
		writeJSONError(res, "Unable to store note")
	},
}
