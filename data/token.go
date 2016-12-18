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

package data

import (
	"errors"
	"time"

	"gopkg.in/hlandau/passlib.v1"
)

const (
	// TokenTypeMaintenace defines maintenance tokens
	TokenTypeMaintenace = 1
	// TokenTypeAccess defines access tokens
	TokenTypeAccess = 2
)

// TokenInterface defines Token
type TokenInterface interface {
	Activate() (Token, error)
	Deactivate() (Token, error)
	IsSecure() bool
	Matches(raw string) bool
	Raw() string
	Remove() error
	Store() (Token, error)
}

// Token implements TokenInterface
type Token struct {
	ID      int       `db:"id"`
	Account int       `db:"account"`
	Text    string    `db:"text"`
	Created time.Time `db:"created"`
	Type    int       `db:"type"`
	Active  bool      `db:"active"`
	raw     string
}

// TokenQueries has all queries for Token
var TokenQueries = map[string]string{
	"tokenAdd": `
		insert into token (account, text, type, active)
		values($1, $2, $3, $4)
		RETURNING id
	`,
	"tokenUpdate": `
		UPDATE token SET text = $2, active = $3
		WHERE id = $1
	`,
	"tokenRemove": `
		delete FROM token WHERE id = $1
	`,
	"tokenGetAllByAccount": `
		SELECT id, account, text, created, type, active FROM token WHERE account = $1
	`,
	"tokenGetAllByAccountAndType": `
		SELECT id, account, text, created, type, active FROM token WHERE account = $1 AND type = $2
	`,
	"tokenGetByID": `
		SELECT id, account, text, created, type, active FROM token WHERE id = $1
	`,
}

// TokenNew creates a new Token
func TokenNew(account int, tokenType int) *Token {
	token := random(32)
	hashed, _ := passlib.Hash(token)

	return &Token{0, account, hashed, time.Now(), tokenType, true, token}
}

// TokenByID retrieves Token by id
func TokenByID(id int) (*Token, error) {
	return tokenByFieldAndValue("tokenGetByID", id)
}

// TokenListByAccountAndType retrieves Token list by Account and type
func TokenListByAccountAndType(account int, tType int) []*Token {
	var list []*Token

	rows, err := pool.Query("tokenGetAllByAccountAndType", account, tType)
	defer rows.Close()

	if err != nil {
		return list
	}

	for rows.Next() {
		token, err := tokenFromResult(rows)

		if err == nil {
			list = append(list, token)
		}
	}

	return list
}

// Activate activates Token and updates the DB
func (t Token) Activate() (*Token, error) {
	if t.Active {
		return &t, nil
	}

	t.Active = true
	return t.Store()
}

// Deactivate activates Token and updates the DB
func (t Token) Deactivate() (*Token, error) {
	if !t.Active {
		return &t, nil
	}

	t.Active = false
	return t.Store()
}

// IsSecure checks Token is secure
func (t Token) IsSecure() bool {
	return t.Raw() == ""
}

// IsStored check if Token is stored in DB
func (t Token) IsStored() bool {
	return t.ID != 0
}

// Matches checks if text matches Token
func (t Token) Matches(raw string) bool {
	_, err := passlib.Verify(raw, t.Text)

	if err == nil {
		return true
	}

	return false
}

// Raw returns Token raw
func (t Token) Raw() string {
	return t.raw
}

// Refresh Token from DB
func (t Token) Refresh() (*Token, error) {
	return TokenByID(t.ID)
}

// Remove Token
func (t Token) Remove() error {
	_, err := pool.Exec("tokenRemove", t.ID)

	return err
}

// Store writes Token to DB
func (t Token) Store() (*Token, error) {
	if t.IsStored() {
		return t.update()
	}

	return t.create()
}

func (t Token) create() (*Token, error) {
	var tokenID int
	err := pool.QueryRow("tokenAdd", t.Account, t.Text, t.Type, t.Active).Scan(&tokenID)

	if err == nil {
		return TokenByID(tokenID)
	}

	return nil, err
}

func (t Token) update() (*Token, error) {
	_, err := pool.Exec("tokenUpdate", t.ID, t.Text, t.Active)

	if err == nil {
		return t.Refresh()
	}

	return nil, err
}

func tokenFromResult(result interface {
	Scan(...interface{}) (err error)
}) (*Token, error) {
	var tokenID int
	var tokenAccount int
	var tokenText string
	var tokenCreated time.Time
	var tokenType int
	var tokenActive bool

	err := result.Scan(
		&tokenID,
		&tokenAccount,
		&tokenText,
		&tokenCreated,
		&tokenType,
		&tokenActive,
	)

	if err == nil {
		return &Token{tokenID, tokenAccount, tokenText, tokenCreated, tokenType, tokenActive, ""}, nil
	}

	return nil, errors.New("Failed to get token")
}

func tokenByFieldAndValue(query string, value interface{}) (*Token, error) {
	return tokenFromResult(pool.QueryRow(query, value))
}
