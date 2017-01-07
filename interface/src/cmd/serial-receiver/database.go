package main

import (
	"database/sql"
	"time"

	"github.com/cenkalti/backoff"
	_ "github.com/lib/pq"
)

const POSTGRES_SETUP_SCRIPT = `
CREATE TABLE IF NOT EXISTS cards (
	card_id varchar(16) UNIQUE NOT NULL,
	user_id integer NOT NULL
);`

var postgres *sql.DB

func setupDatabase() error {
	_, err := postgres.Exec(POSTGRES_SETUP_SCRIPT)
	return err
}

func openPostgres(address string) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute

	return backoff.Retry(func() (err error) {
		postgres, err = sql.Open("postgres", address)
		if err != nil {
			return err
		}

		return postgres.Ping()
	}, b)
}
