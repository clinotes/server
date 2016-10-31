/**
 * clinot.es server
 * Copyright (C) 2016 Sebastian Müller
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.

 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package route

import (
	"encoding/json"
	"errors"
	"net/http"

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

// APIResponseData is
type APIResponseData struct {
	Data  interface{}
	Error bool `json:"error"`
	Done  bool `json:"done"`
}

func writeJSONResponse(res http.ResponseWriter) {
	res.Write([]byte(`{"error": false, "done": true}`))
}

func writeJSONResponseData(res http.ResponseWriter, data interface{}) {
	slcB, _ := json.Marshal(APIResponseData{data, false, !false})
	res.Write([]byte(string(slcB)))
}

func writeJSONError(res http.ResponseWriter, text string) error {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(`{"error", true, "text": "` + text + `"}`))

	return errors.New(text)
}

// Routes returns available routes
func Routes(p *pgx.ConnPool, pm *postmark.Client) []Route {
	pool = p
	pmark = pm

	return []Route{
		APIRouteAdd,
		APIRouteAuth,
		APIRouteAccountCreate,
		APIRouteAccountVerify,
		APIRouteTokenCreate,
		APIRouteSubscribe,
		APIRouteAccount,
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
