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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNote(t *testing.T) {
	acc := AccountNew("mail@example.com")
	user, err := acc.Store()

	assert.Nil(t, err)

	note := NoteNew(user.ID, "example content")

	assert.False(t, note.IsStored())
	assert.Equal(t, "example content", note.Text)
	assert.Equal(t, user.ID, note.Account)

	note, err = note.Store()

	if assert.Nil(t, err) {
		assert.True(t, note.IsStored())
	}

	user.Remove()
}

func TestNoteLimit(t *testing.T) {
	acc := AccountNew("mail@example.com")
	user, err := acc.Store()

	assert.Nil(t, err)

	note := NoteNew(user.ID, "This is a note! This is a note! This is a note! This is a note! This is a note! This is a note! This is a note!")

	assert.False(t, note.IsStored())
	assert.Equal(t, user.ID, note.Account)

	note, err = note.Store()

	assert.NotNil(t, err)

	user.Remove()
}

func TestNoteList(t *testing.T) {
	acc := AccountNew("mail@example.com")
	user, err := acc.Store()

	assert.Nil(t, err)

	note := NoteNew(user.ID, "This is a note!")
	note, err = note.Store()

	assert.Nil(t, err)

	note2 := NoteNew(user.ID, "This is a second note!")
	note2, err = note2.Store()

	assert.Nil(t, err)

	list, err := NoteListByAccount(user.ID)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))

	user.Remove()
}
