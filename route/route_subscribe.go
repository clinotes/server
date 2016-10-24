package route

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
	"github.com/stripe/stripe-go/token"
)

// APIRequestStructSubscribe is
type APIRequestStructSubscribe struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Numer   string `json:"number"`
	Expire  string `json:"expire"`
	CVC     string `json:"cvc"`
}

// APIRouterSubscribe is
var APIRouterSubscribe = Route{
	"/subscribe",
	func(res http.ResponseWriter, req *http.Request) {
		// Parse JSON request
		var data APIRequestStructSubscribe
		if ensureJSONPayload(req, res, &data) != nil {
			return
		}

		// Lookup account ID
		fmt.Println(data)
		accountID, err := accountIDByAddress(data.Address)
		if err != nil {
			fmt.Println("error account id")
			writeJSONError(res, "Unknown account address")
			return
		}

		stripe.Key = os.Getenv("STRIPE_API_KEY")
		expire := strings.Split(data.Expire, "/")
		t, err := token.New(&stripe.TokenParams{
			Card: &stripe.CardParams{
				Number: data.Numer,
				Month:  expire[0],
				Year:   expire[1],
				CVC:    data.CVC,
			},
		})

		if err != nil {
			writeJSONError(res, "Invalid card information")
			return
		}

		customerParams := &stripe.CustomerParams{
			Desc: "Customer for " + data.Address,
		}
		customerParams.SetSource(t.ID)
		c, err := customer.New(customerParams)

		if err != nil {
			writeJSONError(res, "Invalid account information")
			return
		}

		s, err := sub.New(&stripe.SubParams{
			Customer: c.ID,
			Plan:     "blue",
		})

		if err != nil {
			writeJSONError(res, "Invalid account information")
			return
		}

		if _, err = pool.Exec("setSubscription", accountID, s.ID); err == nil {
			writeJSONResponse(res)
			return
		}
		fmt.Println(err)

		writeJSONError(res, "Unable to subscribe")
	},
}
