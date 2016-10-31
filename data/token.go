package data

import (
	"errors"
	"math/rand"
	"time"

	"gopkg.in/hlandau/passlib.v1"
)

const (
	// TokenTypeMaintenace is
	TokenTypeMaintenace = 1
	// TokenTypeAccess is
	TokenTypeAccess = 2
)

// TokenQueries has all queries for account access
var TokenQueries = map[string]string{
	"tokenAdd": `
		insert into token (account, text, type, active)
		values($1, $2, $3, $4)
		RETURNING id
	`,
	"tokenUpdate": `
		UPDATE token SET text = $2, active = $3
		WHERE id = $1
	`,
	"tokenRemove": `
		delete FROM token WHERE id = $1
	`,
	"tokenGetAllByAccount": `
		SELECT id, account, text, created, type, active FROM token WHERE account = $1
	`,
	"tokenGetAllByAccountAndType": `
		SELECT id, account, text, created, type, active FROM token WHERE account = $1 AND type = $2
	`,
	"tokenGetByID": `
		SELECT id, account, text, created, type, active FROM token WHERE id = $1
	`,
}

// TokenInterface is
type TokenInterface interface {
	ID() int
	Account() int
	Text() string
	CreatedOn() time.Time
	Type() int
	IsActive() bool
	Activate() (TokenInterface, error)
	Deactivate() (TokenInterface, error)
	Store() (TokenInterface, error)
	IsSecure() bool
	Remove() error
	Raw() string
	Matches(raw string) bool
}

// Token is
type Token struct {
	id        int
	account   int
	text      string
	created   time.Time
	tokenType int
	active    bool
	raw       string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// Get returns `n` random characters
func random(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// TokenNew generates a new token for account
func TokenNew(account int, tokenType int) TokenInterface {
	token := random(32)
	hashed, _ := passlib.Hash(token)

	return &Token{0, account, hashed, time.Now(), tokenType, true, token}
}

// IsSecure checks if the raw unhashed token is available
func (t Token) IsSecure() bool {
	return t.Raw() == ""
}

// Matches returns true if the passed raw string matches the hashed token
func (t Token) Matches(raw string) bool {
	_, err := passlib.Verify(raw, t.Text())

	if err == nil {
		return true
	}

	return false
}

// Activate token
func (t Token) Activate() (TokenInterface, error) {
	if t.IsActive() {
		return t, nil
	}

	t.active = true
	return t.Store()
}

// Deactivate token
func (t Token) Deactivate() (TokenInterface, error) {
	if !t.IsActive() {
		return t, nil
	}

	t.active = false
	return t.Store()
}

// Remove a token
func (t Token) Remove() error {
	_, err := pool.Exec("tokenRemove", t.ID())

	return err
}

// Store writes the account to the database
func (t Token) Store() (TokenInterface, error) {
	if t.IsStored() {
		return t.update()
	}

	return t.create()
}

func (t Token) create() (TokenInterface, error) {
	var tokenID int
	err := pool.QueryRow("tokenAdd", t.Account(), t.Text(), t.Type(), t.IsActive()).Scan(&tokenID)

	if err == nil {
		return TokenByID(tokenID)
	}

	return nil, err
}

func (t Token) update() (TokenInterface, error) {
	_, err := pool.Exec("tokenUpdate", t.ID(), t.Text(), t.IsActive())

	if err == nil {
		return t.Refresh()
	}

	return nil, err
}

// Raw returns the unhashed string if available
func (t Token) Raw() string {
	return t.raw
}

// IsStored returns true if account is from database
func (t Token) IsStored() bool {
	return t.ID() != 0
}

// Refresh loads gets the token again from DB
func (t Token) Refresh() (*Token, error) {
	return TokenByID(t.ID())
}

// TokenByID returns an Token
func TokenByID(id int) (*Token, error) {
	return tokenByFieldAndValue("tokenGetByID", id)
}

// TokenListByAccountAndType is
func TokenListByAccountAndType(account int, tType int) []TokenInterface {
	var list []TokenInterface

	rows, err := pool.Query("tokenGetAllByAccountAndType", account, tType)
	defer rows.Close()

	if err != nil {
		return list
	}

	for rows.Next() {
		token, err := tokenFromResult(rows)

		if err == nil {
			list = append(list, token)
		}
	}

	return list
}

func tokenFromResult(result interface {
	Scan(...interface{}) (err error)
}) (*Token, error) {
	var tokenID int
	var tokenAccount int
	var tokenText string
	var tokenCreated time.Time
	var tokenType int
	var tokenActive bool

	err := result.Scan(
		&tokenID,
		&tokenAccount,
		&tokenText,
		&tokenCreated,
		&tokenType,
		&tokenActive,
	)

	if err == nil {
		return &Token{tokenID, tokenAccount, tokenText, tokenCreated, tokenType, tokenActive, ""}, nil
	}

	return nil, errors.New("Failed to get token")
}

func tokenByFieldAndValue(query string, value interface{}) (*Token, error) {
	return tokenFromResult(pool.QueryRow(query, value))
}

// ID returns the token ID
func (t Token) ID() int {
	return t.id
}

// IsActive is
func (t Token) IsActive() bool {
	return t.active
}

// CreatedOn is
func (t Token) CreatedOn() time.Time {
	return t.created
}

// Account is
func (t Token) Account() int {
	return t.account
}

// Type is
func (t Token) Type() int {
	return t.tokenType
}

// Text is
func (t Token) Text() string {
	return t.text
}
