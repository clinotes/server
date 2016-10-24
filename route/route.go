package route

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx"
	"github.com/keighl/postmark"
)

var (
	pool  *pgx.ConnPool
	pmark *postmark.Client
)

// Route is a route
type Route struct {
	URL     string
	Handler func(res http.ResponseWriter, req *http.Request)
}

func writeJSONResponse(res http.ResponseWriter) {
	res.Write([]byte(`{"error": false, "done": true}`))
}

func writeJSONResponseData(res http.ResponseWriter, data interface{}) {
	slcB, _ := json.Marshal(APIResponseData{data, false})
	res.Write([]byte(string(slcB)))
}

func writeJSONError(res http.ResponseWriter, text string) error {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(`{"error", true, "text": "` + text + `"}`))

	return errors.New(text)
}

func tokenByUnverifiedAddress(address string) (string, error) {
	var token string

	err := pool.QueryRow("getUnverifiedUser", address).Scan(&token)
	if err != nil {
		return "", errors.New("Failed to retrieve verification token")
	}
	return token, nil
}

// APIResponseData is
type APIResponseData struct {
	Data  interface{}
	Error bool `json:"error"`
}

// APIResponseDataMe is
type APIResponseDataMe struct {
	Address      string
	Created      time.Time
	Notes        int
	Token        int
	Subscription bool
}

func accountByID(accountID int) (*APIResponseDataMe, error) {
	var accountAddress string
	var accountCreated time.Time
	var accountPaid bool

	err := pool.QueryRow("getAccount", accountID).Scan(&accountAddress, &accountCreated, &accountPaid)
	if err != nil {
		return nil, err
	}

	numberToken, err := countTokenByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	numberNotes, err := countNotesByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	return &APIResponseDataMe{accountAddress, accountCreated, numberNotes, numberToken, accountPaid}, nil
}

func countNotesByAccountID(accountID int) (int, error) {
	var count int

	err := pool.QueryRow("countNotes", accountID).Scan(&count)
	if err != nil {
		return 0, errors.New("Failed to count notes for account ID")
	}

	return count, nil
}

func countTokenByAccountID(accountID int) (int, error) {
	var count int

	err := pool.QueryRow("countToken", accountID).Scan(&count)
	if err != nil {
		return 0, errors.New("Failed to count token for account ID")
	}

	return count, nil
}

func accountIDByAddress(address string) (int, error) {
	var accountID int

	err := pool.QueryRow("getUser", address).Scan(&accountID)
	if err != nil {
		return 0, errors.New("Failed to retrieve account ID")
	}

	return accountID, nil
}

// List returns available routes
func List(p *pgx.ConnPool, pm *postmark.Client) []Route {
	pool = p
	pmark = pm

	return []Route{
		APIRouterMe,
		APIRouterAdd,
		APIRouterAuth,
		APIRouteCreateToken,
		APIRouteCreateUser,
		APIRouterVerifyUser,
		APIRouterSubscribe,
	}
}

func ensureJSONPayload(req *http.Request, res http.ResponseWriter, data interface{}) error {
	// Set JSON response header
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Decode body
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	defer req.Body.Close()

	// Respond with BadRequest status
	if err != nil {
		return writeJSONError(res, "Invalid JSON data")
	}

	// Return nil if everything is fine
	return nil
}
