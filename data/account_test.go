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

func TestAccount(t *testing.T) {
	acc := AccountNew("lorem@example.com")

	assert.Equal(t, acc.Address, "lorem@example.com")
	assert.Equal(t, acc.ID, 0)
	assert.False(t, acc.Verified)
	assert.False(t, acc.IsStored())
	assert.NotNil(t, acc.Created)
	assert.False(t, acc.HasSubscription())
	assert.Nil(t, acc.GetSubscription())

	acc, err := acc.Store()

	if assert.Nil(t, err) {
		acc2, err2 := AccountByID(acc.ID)

		assert.Nil(t, err2)
		assert.False(t, acc.Verified)
		assert.Equal(t, acc.ID, acc2.ID)
		assert.Equal(t, acc.Address, acc2.Address)
		assert.Equal(t, acc.Created, acc2.Created)
		assert.Equal(t, acc.Verified, acc2.Verified)

		assert.True(t, acc.IsStored())
		assert.NotEqual(t, acc.ID, 0)
		assert.False(t, acc.Verified)

		acc, err = acc.Verify()
		if assert.Nil(t, err) {
			assert.True(t, acc.Verified)

			assert.NotEqual(t, acc.Verified, acc2.Verified)

			acc3, err3 := acc2.Refresh()

			assert.Nil(t, err3)
			assert.Equal(t, acc.Verified, acc3.Verified)
		}

		acc, err = acc.Store()
		assert.Nil(t, err)

		assert.Equal(t, acc.Address, "lorem@example.com")
		assert.True(t, acc.Verified)

		assert.Equal(t, 0, len(acc.GetTokenList(TokenTypeAccess)))
		assert.Equal(t, 0, len(acc.GetTokenList(TokenTypeMaintenace)))

		tok, err := acc.GetToken("", TokenTypeAccess)

		assert.Nil(t, tok)
		assert.NotNil(t, err)
	}

	acc.Remove()
}
