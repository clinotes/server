package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	acc := AccountNew("lorem@example.com")

	assert.Equal(t, acc.Address(), "lorem@example.com")
	assert.Equal(t, acc.ID(), 0)
	assert.False(t, acc.IsVerified())
	assert.False(t, acc.IsStored())
	assert.NotNil(t, acc.CreatedOn())
	assert.False(t, acc.HasSubscription())
	assert.Nil(t, acc.GetSubscription())

	acc, err := acc.Store()

	if assert.Nil(t, err) {
		acc2, err2 := AccountByID(acc.ID())

		assert.Nil(t, err2)
		assert.False(t, acc.IsVerified())
		assert.Equal(t, acc.ID(), acc2.ID())
		assert.Equal(t, acc.Address(), acc2.Address())
		assert.Equal(t, acc.CreatedOn(), acc2.CreatedOn())
		assert.Equal(t, acc.IsVerified(), acc2.IsVerified())

		assert.True(t, acc.IsStored())
		assert.NotEqual(t, acc.ID(), 0)
		assert.False(t, acc.IsVerified())

		acc, err = acc.Verify()
		if assert.Nil(t, err) {
			assert.True(t, acc.IsVerified())

			assert.NotEqual(t, acc.IsVerified(), acc2.IsVerified())

			acc3, err3 := acc2.Refresh()

			assert.Nil(t, err3)
			assert.Equal(t, acc.IsVerified(), acc3.IsVerified())
		}

		acc, err = acc.Store()
		assert.Nil(t, err)

		assert.Equal(t, acc.Address(), "lorem@example.com")
		assert.True(t, acc.IsVerified())

		assert.Equal(t, 0, len(acc.GetTokenList(TokenTypeAccess)))
		assert.Equal(t, 0, len(acc.GetTokenList(TokenTypeMaintenace)))
	}

	acc.Remove()
}
