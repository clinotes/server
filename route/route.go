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

	"github.com/keighl/postmark"
)

// Handler is
type Handler func(http.ResponseWriter, *http.Request) (interface{}, error)

// APIResponseData is
type apiResponseData struct {
	Data  interface{}
	Error bool `json:"error"`
	Done  bool `json:"done"`
}

// APIResponseSuccess is a successful API response
type apiResponseSuccess struct {
	Error bool `json:"error"`
	Done  bool `json:"done"`
}

// APIResponseError is an API error
type apiResponseError struct {
	Error bool   `json:"error"`
	Text  string `json:"text"`
}

// Route is a route
type Route struct {
	URL     string
	Handler Handler
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set JSON response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Prepare response object
	var response interface{}

	// Check for error in route handler
	data, err := handler(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = apiResponseError{true, err.Error()}
	} else {
		if data != nil {
			response = apiResponseData{data, false, !false}
		} else {
			response = apiResponseSuccess{false, true}
		}
	}

	// Write response
	text, _ := json.Marshal(response)
	w.Write([]byte(string(text)))
}

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

func checkJSONBody(req *http.Request, res http.ResponseWriter, data interface{}) error {
	// Decode body
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	defer req.Body.Close()

	// Respond with BadRequest status
	if err != nil {
		return errors.New("Invalid JSON data")
	}

	return nil
}

func sendTokenWithTemplate(to string, token string, template int64) (postmark.EmailResponse, error) {
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
