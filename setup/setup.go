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
