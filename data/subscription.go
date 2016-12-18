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

import "time"

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

// SubscriptionNew creates a new Subscription
func SubscriptionNew(account int, stripeid string) *Subscription {
	return &Subscription{0, account, time.Now(), stripeid, false}
}

// SubscriptionByID retrieves Subscription by id
func SubscriptionByID(id int) (*Subscription, error) {
	var sub Subscription

	err := db.Get(&sub, "SELECT id, account, created, stripeid, active FROM subscription WHERE id = $1", id)

	return &sub, err
}

// SubscriptionByAccountID retrieves Subscription by Account id
func SubscriptionByAccountID(id int) (*Subscription, error) {
	var sub Subscription

	err := db.Select(&sub, `SELECT id, account, created, stripeid, active
		FROM subscription WHERE account = $1 AND active = TRUE`, id)

	return &sub, err
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
	var id int
	rows, err := db.Query(`
		insert into subscription (account, stripeid)
		values($1, $2)
		RETURNING id
	`, s.Account, s.StripeID)

	if err != nil {
		return nil, err
	}

	rows.Next()
	rows.Scan(&id)

	return SubscriptionByID(id)
}

func (s Subscription) update() (*Subscription, error) {
	_, err := db.Query(`UPDATE subscription SET active = $2
		WHERE id = $1`, s.ID, s.Active)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
