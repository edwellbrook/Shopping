package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	_ "github.com/lib/pq"
)

var postgres *sql.DB

func openPostgres() (err error) {
	postgres, err = sql.Open("postgres", config.Postgres)
	if err != nil {
		log.Println(err)
	} else {
		err = postgres.Ping()
	}
	return err
}

func setupDatabase() {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute

	err := backoff.Retry(openPostgres, b)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %s\n", err)
	}

	script := `
CREATE TABLE IF NOT EXISTS cards (
	card_id varchar(16) not null,
	user_id integer not null
);`

	if _, err = postgres.Exec(script); err != nil {
		log.Fatal(err)
	}
}
