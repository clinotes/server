package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscription(t *testing.T) {
	acc := AccountNew("lorem@example.com")
	user, err := acc.Store()

	assert.Nil(t, err)

	sub := SubscriptionNew(user.ID(), "test")

	assert.Equal(t, 0, sub.ID())
	assert.Equal(t, user.ID(), sub.Account())
	assert.Equal(t, "test", sub.StripeID())
	assert.False(t, sub.IsActive())
	assert.False(t, sub.IsStored())

	sub, err = sub.Store()

	if assert.Nil(t, err) {
		assert.NotEqual(t, 0, sub.ID())
		assert.False(t, sub.IsActive())
		assert.True(t, sub.IsStored())

		sub, err = sub.Activate()
		if assert.Nil(t, err) {
			assert.True(t, sub.IsActive())
		}

		sub, err = sub.Activate()
		if assert.Nil(t, err) {
			assert.True(t, sub.IsActive())
		}

		sub, err = sub.Deactivate()
		if assert.Nil(t, err) {
			assert.False(t, sub.IsActive())
		}

		sub, err = sub.Deactivate()
		if assert.Nil(t, err) {
			assert.False(t, sub.IsActive())
		}
	}

	user.Remove()
}
