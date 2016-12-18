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

// SubscriptionInterface defines Subscription
type SubscriptionInterface interface {
	Activate() (*Subscription, error)
	Deactivate() (*Subscription, error)
	IsStored() bool
	Refresh() (*Subscription, error)
	Store() (*Subscription, error)

	create() (SubscriptionInterface, error)
	update() (SubscriptionInterface, error)
}

// Subscription implements SubscriptionInterface
type Subscription struct {
	ID       int       `db:"id"`
	Account  int       `db:"account"`
	Created  time.Time `db:"created"`
	StripeID string    `db:"stripeid"`
	Active   bool      `db:"active"`
}

// SubscriptionQueries has all queries for Subscription
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

// SubscriptionNew creates a new Subscription
func SubscriptionNew(account int, stripeid string) *Subscription {
	return &Subscription{0, account, time.Now(), stripeid, false}
}

// SubscriptionByID retrieves Subscription by id
func SubscriptionByID(id int) (*Subscription, error) {
	return subscriptionByFieldAndValue("subscriptionGetByID", id)
}

// SubscriptionByAccountID retrieves Subscription by Account id
func SubscriptionByAccountID(id int) (*Subscription, error) {
	return subscriptionByFieldAndValue("subscriptionGetByAccountID", id)
}

// Activate activates Subscripiton and updates the DB
func (s Subscription) Activate() (*Subscription, error) {
	if s.Active {
		return &s, nil
	}

	s.Active = true
	return s.Store()
}

// Deactivate deactivates Subscrition and updates the DB
func (s Subscription) Deactivate() (*Subscription, error) {
	if !s.Active {
		return &s, nil
	}

	s.Active = false
	return s.Store()
}

// IsStored checks if Subscription is stored in DB
func (s Subscription) IsStored() bool {
	return s.ID != 0
}

// Refresh Subscription from DB
func (s Subscription) Refresh() (*Subscription, error) {
	return SubscriptionByID(s.ID)
}

// Store writes Subscription to DB
func (s Subscription) Store() (*Subscription, error) {
	if s.IsStored() {
		return s.update()
	}

	return s.create()
}

func (s Subscription) create() (*Subscription, error) {
	var subscriptionID int
	err := pool.QueryRow("subscriptionAdd", s.Account, s.StripeID).Scan(&subscriptionID)

	if err == nil {
		return SubscriptionByID(subscriptionID)
	}

	return nil, err
}

func (s Subscription) update() (*Subscription, error) {
	_, err := pool.Exec("subscriptionUpdate", s.ID, s.Active)

	if err == nil {
		return s.Refresh()
	}

	return nil, err
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
