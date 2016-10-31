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

package setup

import "github.com/jackc/pgx"

// Run creates the database structure if needed
func Run(pool *pgx.Conn) {
	pool.Exec(`
		CREATE TABLE account(
	    id serial primary key,
	    address TEXT NOT NULL,
	    created TIMESTAMP DEFAULT now() NOT NULL,
	    verified BOOLEAN DEFAULT false NOT NULL
		);

		CREATE TABLE subscription(
	    id serial primary key,
	    account INTEGER NOT NULL,
	    created TIMESTAMP DEFAULT now() NOT NULL,
	    stripeid TEXT NOT NULL,
	    active BOOLEAN DEFAULT false
		);

		CREATE TABLE token(
	    id serial primary key,
	    account INTEGER NOT NULL,
	    text TEXT NOT NULL,
	    created TIMESTAMP DEFAULT now() NOT NULL,
	    type INTEGER DEFAULT 1 NOT NULL,
	    active BOOLEAN DEFAULT true
		);

		CREATE UNIQUE INDEX account_id_uindex ON account (id);
		CREATE UNIQUE INDEX account_address_uindex ON account (address);

		ALTER TABLE subscription ADD FOREIGN KEY (account) REFERENCES account (id) on delete cascade;
		CREATE UNIQUE INDEX subscription_id_uindex ON subscription (id);
		CREATE UNIQUE INDEX "subscription_stripeID_uindex" ON subscription (stripeid);

		ALTER TABLE token ADD FOREIGN KEY (account) REFERENCES account (id) on delete cascade;
		CREATE UNIQUE INDEX token_id_uindex ON token (id);
	`)
}
