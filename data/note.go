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

// NoteInterface defines Note
type NoteInterface interface {
	Account() int
	CreatedOn() time.Time
	ID() int
	IsStored() bool
	Store() (NoteInterface, error)
	Text() string

	create() (NoteInterface, error)
	update() (NoteInterface, error)
}

// Note implements NoteInterface
type Note struct {
	id      int
	account int
	text    string
	created time.Time
}

// NoteQueries has all queries for Note
var NoteQueries = map[string]string{
	"noteAdd": `
		insert into note (account, text)
		values($1, $2)
		RETURNING id
	`,
	"noteUpdate": `
		UPDATE note SET text = $2
		WHERE id = $1
	`,
	"noteGetByID": `
		SELECT id, account, text, created FROM note WHERE id = $1
	`,
}

// NoteNew creates a new Note
func NoteNew(account int, text string) NoteInterface {
	return Note{0, account, text, time.Now()}
}

// NoteByID retrieves Note by id
func NoteByID(id int) (*Note, error) {
	return noteByFieldAndValue("noteGetByID", id)
}

// Account retrieves Note account
func (n Note) Account() int {
	return n.account
}

// CreatedOn returns Note create date
func (n Note) CreatedOn() time.Time {
	return n.created
}

// ID returns Account id
func (n Note) ID() int {
	return n.id
}

// IsStored checks if Note is stored in DB
func (n Note) IsStored() bool {
	return n.ID() != 0
}

// Refresh Note from DB
func (n Note) Refresh() (*Note, error) {
	return NoteByID(n.ID())
}

// Store writes Notes to DB
func (n Note) Store() (NoteInterface, error) {
	if len(n.Text()) > 100 {
		return nil, errors.New("Note must not be longer than 100 characters")
	}
	if n.IsStored() {
		return n.update()
	}

	return n.create()
}

// Text returns Note text
func (n Note) Text() string {
	return n.text
}

func (n Note) create() (NoteInterface, error) {
	var noteID int
	err := pool.QueryRow("noteAdd", n.Account(), n.Text()).Scan(&noteID)

	if err == nil {
		return NoteByID(noteID)
	}

	return nil, err

}

func (n Note) update() (NoteInterface, error) {
	_, err := pool.Exec("noteUpdate", n.ID(), n.Text())

	if err == nil {
		return n.Refresh()
	}

	return nil, err
}

func noteFromResult(result interface {
	Scan(...interface{}) (err error)
}) (*Note, error) {
	var noteID int
	var noteAccount int
	var noteText string
	var noteCreated time.Time

	err := result.Scan(
		&noteID,
		&noteAccount,
		&noteText,
		&noteCreated,
	)

	if err == nil {
		return &Note{noteID, noteAccount, noteText, noteCreated}, nil
	}

	return nil, errors.New("Failed to get token")
}

func noteByFieldAndValue(query string, value interface{}) (*Note, error) {
	return noteFromResult(pool.QueryRow(query, value))
}
