package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Route is a route
type Route struct {
	URL     string
	Handler func(res http.ResponseWriter, req *http.Request)
}

func writeJSONResponse(res http.ResponseWriter) {
	res.Write([]byte(`{"error": false}`))
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

func accountIDByAddress(address string) (int, error) {
	var accountID int

	err := pool.QueryRow("getUser", address).Scan(&accountID)
	if err != nil {
		return 0, errors.New("Failed to retrieve account ID")
	}

	return accountID, nil
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
