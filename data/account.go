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
	Address() string
	CreatedOn() time.Time
	GetSubscription() SubscriptionInterface
	GetToken(t string, tokenType int) (*Token, error)
	GetTokenList(tokenType int) []*Token
	HasSubscription() bool
	ID() int
	IsStored() bool
	IsVerified() bool
	Refresh() (*Account, error)
	Remove() error
	Store() (AccountInterface, error)
	Verify() (AccountInterface, error)

	create() (AccountInterface, error)
	update() (AccountInterface, error)
}

// Account implements AccountInterface
type Account struct {
	id       int
	address  string
	created  time.Time
	verified bool
}

// AccountQueries has all queries for Account
var AccountQueries = map[string]string{
	"accountAdd": `
		insert into account (address)
		values($1)
	`,
	"accountRemove": `
		delete FROM account WHERE id = $1
	`,
	"accountGetByAddress": `
		SELECT id, address, created, verified FROM account WHERE address = $1
	`,
	"accountGetByID": `
		SELECT id, address, created, verified FROM account WHERE id = $1
	`,
	"accountUpdate": `
		UPDATE account SET verified = $2
		WHERE id = $1
	`,
}

// AccountNew creates a new account
func AccountNew(address string) AccountInterface {
	return &Account{0, address, time.Now(), false}
}

// AccountByAddress retrieves Account by address
func AccountByAddress(address string) (AccountInterface, error) {
	return accountByFieldAndValue("accountGetByAddress", address)
}

// AccountByID retrieves Account by id
func AccountByID(id int) (*Account, error) {
	return accountByFieldAndValue("accountGetByID", id)
}

// Address returns Account address
func (a Account) Address() string {
	return a.address
}

// CreatedOn returns Account create date
func (a Account) CreatedOn() time.Time {
	return a.created
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
	return TokenListByAccountAndType(a.ID(), tokenType)
}

// GetSubscription retrieves Account Subscription
func (a Account) GetSubscription() SubscriptionInterface {
	sub, err := SubscriptionByAccountID(a.ID())

	if err == nil {
		return sub
	}

	return nil
}

// HasSubscription checks if Account has a Subscription
func (a Account) HasSubscription() bool {
	return a.GetSubscription() != nil
}

// ID returns Account id
func (a Account) ID() int {
	return a.id
}

// IsStored checks if Account is stored in DB
func (a Account) IsStored() bool {
	return a.ID() != 0
}

// IsVerified checks if Account is verified
func (a Account) IsVerified() bool {
	return a.verified
}

// Refresh Account from DB
func (a Account) Refresh() (*Account, error) {
	return AccountByID(a.ID())
}

// Remove Account
func (a Account) Remove() error {
	_, err := pool.Exec("accountRemove", a.ID())

	return err
}

// Store writes Account to DB
func (a Account) Store() (AccountInterface, error) {
	if a.IsStored() {
		return a.update()
	}

	return a.create()
}

// Verify verifies Account and updates the DB
func (a Account) Verify() (AccountInterface, error) {
	_, err := pool.Exec("accountUpdate", a.ID(), true)

	if err != nil {
		return nil, err
	}

	return AccountByID(a.ID())
}

func (a Account) create() (AccountInterface, error) {
	_, err := pool.Exec("accountAdd", a.Address())

	if err == nil {
		return AccountByAddress(a.Address())
	}

	return nil, err
}

func (a Account) update() (AccountInterface, error) {
	_, err := pool.Exec("accountUpdate", a.ID(), a.IsVerified())

	if err != nil {
		return nil, err
	}

	return AccountByID(a.ID())
}

func accountFromResult(result interface {
	Scan(...interface{}) (err error)
}) (*Account, error) {
	var accountID int
	var accountAddress string
	var accountCreated time.Time
	var accountVerified bool

	err := result.Scan(
		&accountID,
		&accountAddress,
		&accountCreated,
		&accountVerified,
	)

	if err == nil {
		return &Account{accountID, accountAddress, accountCreated, accountVerified}, nil
	}

	return nil, errors.New("Failed to get subscription")
}

func accountByFieldAndValue(query string, value interface{}) (*Account, error) {
	return accountFromResult(pool.QueryRow(query, value))
}
