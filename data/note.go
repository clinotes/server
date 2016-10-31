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

// NoteInterface is
type NoteInterface interface {
	Text() string
	CreatedOn() time.Time
}

// Note is
type Note struct {
	text    string
	created time.Time
}

// NoteNew returns
func NoteNew(note string) NoteInterface {
	return Note{note, time.Now()}
}

// Text is
func (n Note) Text() string {
	return n.text
}

// CreatedOn is
func (n Note) CreatedOn() time.Time {
	return n.created
}
