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

// TokenNew creates a new Token
func TokenNew(account int, tokenType int) *Token {
	token := random(32)
	hashed, _ := passlib.Hash(token)

	return &Token{0, account, hashed, time.Now(), tokenType, true, token}
}

// TokenByID retrieves Token by id
func TokenByID(id int) (*Token, error) {
	var token Token

	err := db.Get(&token, "SELECT id, account, text, created, type, active FROM token WHERE id = $1", id)

	return &token, err
}

// TokenListByAccountAndType retrieves Token list by Account and type
func TokenListByAccountAndType(account int, tType int) []*Token {
	var list []*Token

	db.Select(&list, `SELECT id, account, text, created, type, active
		FROM token WHERE account = $1 AND type = $2`, account, tType)

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

// Remove Token
func (t Token) Remove() error {
	_, err := db.Exec("delete FROM token WHERE id = $1", t.ID)

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
	var id int
	rows, err := db.Query(`
		insert into token (account, text, type, active)
		values($1, $2, $3, $4)
		RETURNING id
	`, t.Account, t.Text, t.Type, t.Active)

	if err != nil {
		return nil, err
	}

	rows.Next()
	rows.Scan(&id)

	return TokenByID(id)
}

func (t Token) update() (*Token, error) {
	_, err := db.Query(`UPDATE token SET text = $2, active = $3
		WHERE id = $1`, t.ID, t.Text, t.Active)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
