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
	"fmt"
	"net/http"

	"github.com/keighl/postmark"
)

var (
	conf      Configuration
	pmark     *postmark.Client
	templates map[string]int64
)

// Configuration stores need variables
type Configuration struct {
	TemplateWelcome int64
	TemplateConfirm int64
	TemplateToken   int64

	PostmarkToken   string
	PostmarkFrom    string
	PostmarkReplyTo string
}

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

// SetTemplate sets template id for a key
func SetTemplate(key string, id int64) {
	templates[key] = id
}

// GetTemplate gets template id for a key
func GetTemplate(key string) int64 {
	return templates[key]
}

// Routes returns available routes
func Routes(config Configuration) []Route {
	conf = config
	pmark = postmark.NewClient(config.PostmarkToken, "")

	return []Route{
		APIRouteAdd,
		APIRouteAuth,
		APIRouteAccountCreate,
		APIRouteAccountVerify,
		APIRouteTokenCreate,
		APIRouteSubscribe,
		APIRouteAccount,
		APIRouteNotes,
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

func sendTokenWithTemplate(to string, token string, template int64) (postmark.EmailResponse, error) {
	fmt.Printf("Using %d for mail\n", template)
	return pmark.SendTemplatedEmail(postmark.TemplatedEmail{
		TemplateId: template,
		TemplateModel: map[string]interface{}{
			"token": token,
		},
		From:    conf.PostmarkFrom,
		To:      to,
		ReplyTo: conf.PostmarkReplyTo,
	})
}
