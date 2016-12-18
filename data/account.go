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
)

// AccountInterface defines Account
type AccountInterface interface {
	GetSubscription() *Subscription
	GetToken(t string, tokenType int) (*Token, error)
	GetTokenList(tokenType int) []*Token
	HasSubscription() bool
	IsStored() bool
	Refresh() (*Account, error)
	Remove() error
	Store() (*Account, error)
	Verify() (*Account, error)

	create() (*Account, error)
	update() (*Account, error)
}

// Account implements AccountInterface
type Account struct {
	ID       int       `db:"id"`
	Address  string    `db:"address"`
	Created  time.Time `db:"created"`
	Verified bool      `db:"verified"`
}

// AccountNew creates a new account
func AccountNew(address string) *Account {
	return &Account{0, address, time.Now(), false}
}

// AccountByAddress retrieves Account by address
func AccountByAddress(address string) (*Account, error) {
	var account Account

	err := db.Get(&account, "SELECT id, address, created, verified FROM account WHERE address = $1", address)

	return &account, err
}

// AccountByID retrieves Account by id
func AccountByID(id int) (*Account, error) {
	var account Account

	err := db.Get(&account, "SELECT id, address, created, verified FROM account WHERE id = $1", id)

	return &account, err
}

// GetToken retrieves Token for Account
func (a Account) GetToken(t string, tokenType int) (*Token, error) {
	token := &Token{}
	found := false

	for _, item := range a.GetTokenList(tokenType) {
		if item.Matches(t) {
			found = true
			token = item
		}
	}

	if found {
		return token, nil
	}

	return nil, errors.New("Token not found")
}

// GetTokenList retrieves all Token for Account
func (a Account) GetTokenList(tokenType int) []*Token {
	return TokenListByAccountAndType(a.ID, tokenType)
}

// GetSubscription retrieves Account Subscription
func (a Account) GetSubscription() *Subscription {
	sub, err := SubscriptionByAccountID(a.ID)

	if err == nil {
		return sub
	}

	return nil
}

// HasSubscription checks if Account has a Subscription
func (a Account) HasSubscription() bool {
	return a.GetSubscription() != nil
}

// IsStored checks if Account is stored in DB
func (a Account) IsStored() bool {
	return a.ID != 0
}

// Refresh Account from DB
func (a Account) Refresh() (*Account, error) {
	return AccountByID(a.ID)
}

// Remove Account
func (a Account) Remove() error {
	_, err := db.Exec("delete FROM account WHERE id = $1", a.ID)

	return err
}

// Store writes Account to DB
func (a Account) Store() (*Account, error) {
	if a.IsStored() {
		return a.update()
	}

	return a.create()
}

// Verify verifies Account and updates the DB
func (a Account) Verify() (*Account, error) {
	a.Verified = true

	return a.update()
}

func (a Account) create() (*Account, error) {
	var id int
	rows, err := db.Query("insert into account (address) values($1)", a.Address)

	if err != nil {
		return nil, err
	}

	rows.Next()
	rows.Scan(&id)

	return AccountByAddress(a.Address)
}

func (a Account) update() (*Account, error) {
	_, err := db.Query(`UPDATE account SET verified = $2
		WHERE id = $1`, a.ID, a.Verified)

	if err != nil {
		return nil, err
	}

	return &a, nil
}
