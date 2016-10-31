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

// SubscriptionInterface defines subscription interactions
type SubscriptionInterface interface {
	ID() int
	Account() int
	CreatedOn() time.Time
	StripeID() string
	IsActive() bool
	Store() (SubscriptionInterface, error)
	Activate() (SubscriptionInterface, error)
	Deactivate() (SubscriptionInterface, error)

	IsStored() bool
}

// Subscription is the general account subscription
type Subscription struct {
	id       int
	account  int
	created  time.Time
	stripeid string
	active   bool
}

// SubscriptionQueries has all queries for account access
var SubscriptionQueries = map[string]string{
	"subscriptionAdd": `
		insert into subscription (account, stripeid)
		values($1, $2)
		RETURNING id
	`,
	"subscriptionUpdate": `
		UPDATE subscription SET active = $2
		WHERE id = $1
	`,
	"subscriptionGetByID": `
		SELECT id, account, created, stripeid, active FROM subscription WHERE id = $1
	`,
	"subscriptionGetByAccountID": `
		SELECT id, account, created, stripeid, active FROM subscription WHERE account = $1 AND active = TRUE
	`,
}

// SubscriptionNew creates a new account
func SubscriptionNew(account int, stripeid string) SubscriptionInterface {
	return &Subscription{0, account, time.Now(), stripeid, false}
}

// Activate token
func (s Subscription) Activate() (SubscriptionInterface, error) {
	if s.IsActive() {
		return s, nil
	}

	s.active = true
	return s.Store()
}

// Deactivate token
func (s Subscription) Deactivate() (SubscriptionInterface, error) {
	if !s.IsActive() {
		return s, nil
	}

	s.active = false
	return s.Store()
}

// ID is
func (s Subscription) ID() int {
	return s.id
}

// Account is
func (s Subscription) Account() int {
	return s.account
}

// CreatedOn returns the date where the subscription was created
func (s Subscription) CreatedOn() time.Time {
	return s.created
}

// StripeID is
func (s Subscription) StripeID() string {
	return s.stripeid
}

// IsActive returns true if the subscription is active
func (s Subscription) IsActive() bool {
	return s.active
}

// IsStored returns true if account is from database
func (s Subscription) IsStored() bool {
	return s.ID() != 0
}

// Store writes the account to the database
func (s Subscription) Store() (SubscriptionInterface, error) {
	if s.IsStored() {
		return s.update()
	}

	return s.create()
}

func (s Subscription) create() (SubscriptionInterface, error) {
	var subscriptionID int
	err := pool.QueryRow("subscriptionAdd", s.Account(), s.StripeID()).Scan(&subscriptionID)

	if err == nil {
		return SubscriptionByID(subscriptionID)
	}

	return nil, err
}

func (s Subscription) update() (SubscriptionInterface, error) {
	_, err := pool.Exec("subscriptionUpdate", s.ID(), s.IsActive())

	if err == nil {
		return s.Refresh()
	}

	return nil, err
}

// Refresh loads gets the token again from DB
func (s Subscription) Refresh() (*Subscription, error) {
	return SubscriptionByID(s.ID())
}

// SubscriptionByID returns an Token
func SubscriptionByID(id int) (*Subscription, error) {
	return subscriptionByFieldAndValue("subscriptionGetByID", id)
}

// SubscriptionByAccountID returns an Token
func SubscriptionByAccountID(id int) (*Subscription, error) {
	return subscriptionByFieldAndValue("subscriptionGetByAccountID", id)
}

func subscriptionFromResult(result interface {
	Scan(...interface{}) (err error)
}) (*Subscription, error) {
	var subscriptionID int
	var subscriptionAccount int
	var subscriptionCreated time.Time
	var subscriptionStripeID string
	var subscriptionActive bool

	err := result.Scan(
		&subscriptionID,
		&subscriptionAccount,
		&subscriptionCreated,
		&subscriptionStripeID,
		&subscriptionActive,
	)

	if err == nil {
		return &Subscription{subscriptionID, subscriptionAccount, subscriptionCreated, subscriptionStripeID, subscriptionActive}, nil
	}

	return nil, errors.New("Failed to get subscription")
}

func subscriptionByFieldAndValue(query string, value interface{}) (*Subscription, error) {
	return subscriptionFromResult(pool.QueryRow(query, value))
}
