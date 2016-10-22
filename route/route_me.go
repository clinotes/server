package route

import "net/http"

// APIRequestStructMe is
type APIRequestStructMe struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// APIRouterMe is
var APIRouterMe = Route{
	"/me",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var data APIRequestStructMe
		if ensureJSONPayload(req, res, &data) != nil {
			return
		}

		// Lookup account ID
		accountID, err := accountIDByAddress(data.Address)
		if err != nil {
			writeJSONError(res, "Unknown account address")
			return
		}

		accountData, err := accountByID(accountID)
		if err != nil {
			writeJSONError(res, "Unknown account address")
			return
		}

		writeJSONResponseData(res, accountData)
	},
}
