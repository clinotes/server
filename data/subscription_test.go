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
