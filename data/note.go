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
