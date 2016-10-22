package route

import (
	"fmt"
	"net/http"
)

// APIRequestStructAdd is
type APIRequestStructAdd struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Note    string `json:"note"`
}

// APIRouterAdd is
var APIRouterAdd = Route{
	"/add",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var data APIRequestStructAdd
		if ensureJSONPayload(req, res, &data) != nil {
			return
		}

		fmt.Printf("Received Note: %s\n", data.Note)
		fmt.Printf(" - %s / %s", data.Address, data.Token)

		// Return error if no token is found
		writeJSONError(res, "Unable to store note")
	},
}
