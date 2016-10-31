package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	acc := AccountNew("mail@example.com")
	user, err := acc.Store()

	token := TokenNew(user.ID(), TokenTypeMaintenace)

	assert.NotEqual(t, "", token.Raw())
	assert.False(t, token.IsSecure())
	assert.True(t, token.Matches(token.Raw()))
	assert.False(t, token.Matches("test"+token.Raw()))
	assert.Equal(t, token.ID(), 0)
	assert.Equal(t, token.Type(), TokenTypeMaintenace)
	assert.NotNil(t, token.Text())
	assert.True(t, token.IsActive())

	token, err = token.Store()

	if assert.Nil(t, err) {
		assert.Equal(t, "", token.Raw())
		assert.True(t, token.IsSecure())
		assert.NotEqual(t, 0, token.ID())
		assert.NotEqual(t, 0, token.CreatedOn())
	}

	token, err = token.Deactivate()
	if assert.Nil(t, err) {
		assert.False(t, token.IsActive())
	}

	token, err = token.Deactivate()
	if assert.Nil(t, err) {
		assert.False(t, token.IsActive())
	}

	token, err = token.Activate()
	if assert.Nil(t, err) {
		assert.True(t, token.IsActive())
	}

	token, err = token.Activate()
	if assert.Nil(t, err) {
		assert.True(t, token.IsActive())
	}

	token.Remove()
	user.Remove()
}

func TestTokenList(t *testing.T) {
	acc := AccountNew("mail@example.com")
	user, err := acc.Store()

	assert.True(t, user.IsStored())
	assert.Nil(t, err)

	token := TokenNew(user.ID(), TokenTypeMaintenace)
	token, err = token.Store()

	assert.Nil(t, err)
	assert.NotNil(t, token.ID())

	token2 := TokenNew(user.ID(), TokenTypeMaintenace)
	token2, err = token2.Store()

	assert.Nil(t, err)
	assert.NotNil(t, token2.ID())

	token3 := TokenNew(user.ID(), TokenTypeAccess)
	token3, err = token3.Store()

	assert.Nil(t, err)
	assert.NotNil(t, token3.ID())

	listMaintenance := TokenListByAccountAndType(user.ID(), TokenTypeMaintenace)
	assert.Equal(t, 2, len(listMaintenance))

	listAccess := TokenListByAccountAndType(user.ID(), TokenTypeAccess)
	assert.Equal(t, 1, len(listAccess))

	user.Remove()
}
