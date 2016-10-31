package data

import (
	"errors"
	"time"
)

// AccountInterface defines account interactions
type AccountInterface interface {
	ID() int
	Address() string
	IsStored() bool
	IsVerified() bool
	CreatedOn() time.Time
	HasSubscription() bool
	GetSubscription() SubscriptionInterface
	Verify() (AccountInterface, error)
	GetToken(t string, tokenType int) (TokenInterface, error)
	GetTokenList(tokenType int) []TokenInterface

	Refresh() (*Account, error)
	create() (AccountInterface, error)
	update() (AccountInterface, error)
	Store() (AccountInterface, error)
	Remove() error
}

// Account is the general user account
type Account struct {
	id       int
	address  string
	created  time.Time
	verified bool
}

// AccountQueries has all queries for account access
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

// GetTokenList returns a list of token
func (a Account) GetTokenList(tokenType int) []TokenInterface {
	return TokenListByAccountAndType(a.ID(), tokenType)
}

// GetToken reads token from user
func (a Account) GetToken(t string, tokenType int) (TokenInterface, error) {
	var token TokenInterface
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

// Remove an account
func (a Account) Remove() error {
	_, err := pool.Exec("accountRemove", a.ID())

	return err
}

// Refresh loads gets the user again from DB
func (a Account) Refresh() (*Account, error) {
	return AccountByID(a.ID())
}

// Verify sets verified to true
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

// Store writes the account to the database
func (a Account) Store() (AccountInterface, error) {
	if a.IsStored() {
		return a.update()
	}

	return a.create()
}

// HasSubscription is
func (a Account) HasSubscription() bool {
	return a.GetSubscription() != nil
}

// IsStored returns true if account is from database
func (a Account) IsStored() bool {
	return a.ID() != 0
}

// GetSubscription is
func (a Account) GetSubscription() SubscriptionInterface {
	sub, err := SubscriptionByAccountID(a.ID())

	if err == nil {
		return sub
	}

	return nil
}

// CreatedOn is
func (a Account) CreatedOn() time.Time {
	return a.created
}

// IsVerified is
func (a Account) IsVerified() bool {
	return a.verified
}

// Address is
func (a Account) Address() string {
	return a.address
}

// ID is
func (a Account) ID() int {
	return a.id
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

// AccountByAddress returns an Account
func AccountByAddress(address string) (AccountInterface, error) {
	return accountByFieldAndValue("accountGetByAddress", address)
}

// AccountByID returns an Account
func AccountByID(id int) (*Account, error) {
	return accountByFieldAndValue("accountGetByID", id)
}
