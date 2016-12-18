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
	IsStored() bool
	Store() (*Note, error)

	create() (*Note, error)
	update() (*Note, error)
}

// Note implements NoteInterface
type Note struct {
	ID      int       `db:"id"`
	Account int       `db:"account"`
	Text    string    `db:"text"`
	Created time.Time `db:"created"`
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
	"noteListGetByAccount": `
		SELECT * FROM (
			SELECT id, account, text, created FROM note WHERE account = $1 ORDER BY id DESC LIMIT 10
		) as list ORDER BY id ASC
	`,
}

// NoteNew creates a new Note
func NoteNew(account int, text string) *Note {
	return &Note{0, account, text, time.Now()}
}

// NoteByID retrieves Note by id
func NoteByID(id int) (*Note, error) {
	return noteByFieldAndValue("noteGetByID", id)
}

// NoteListByAccount retrieves Note by id
func NoteListByAccount(account int) ([]Note, error) {
	return noteListByFieldAndValue("noteListGetByAccount", account)
}

// IsStored checks if Note is stored in DB
func (n Note) IsStored() bool {
	return n.ID != 0
}

// Refresh Note from DB
func (n Note) Refresh() (*Note, error) {
	return NoteByID(n.ID)
}

// Store writes Notes to DB
func (n Note) Store() (*Note, error) {
	if len(n.Text) > 100 {
		return nil, errors.New("Note must not be longer than 100 characters")
	}

	if n.IsStored() {
		return n.update()
	}

	return n.create()
}

func (n Note) create() (*Note, error) {
	var noteID int
	err := pool.QueryRow("noteAdd", n.Account, n.Text).Scan(&noteID)

	if err == nil {
		return NoteByID(noteID)
	}

	return nil, err

}

func (n Note) update() (*Note, error) {
	_, err := pool.Exec("noteUpdate", n.ID, n.Text)

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

func noteListFromResult(result interface {
	Next() bool
	Scan(dest ...interface{}) (err error)
}) ([]Note, error) {
	var noteID int
	var noteAccount int
	var noteText string
	var noteCreated time.Time
	var list []Note

	for result.Next() {
		err := result.Scan(
			&noteID,
			&noteAccount,
			&noteText,
			&noteCreated,
		)

		if err != nil {
			return nil, errors.New("Failed to get notes")
		}

		list = append(list, Note{noteID, noteAccount, noteText, noteCreated})
	}

	return list, nil
}

func noteListByFieldAndValue(query string, value interface{}) ([]Note, error) {
	q, err := pool.Query(query, value)

	if err != nil {
		return nil, errors.New("Failed to get notes")
	}

	return noteListFromResult(q)
}

func noteByFieldAndValue(query string, value interface{}) (*Note, error) {
	return noteFromResult(pool.QueryRow(query, value))
}
