package migrations

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Run(sourceURL string, databaseURL string) {
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatalf("Unable to run migrations: %v\n", err)
	}
	if err := m.Up(); err != nil {
		// exit the app if its not 'first : file does not exist' or 'no change' error
		if (err.Error() != "first : file does not exist") && (err.Error() != "no change") {
			log.Fatalf("Unable to run migrations: %v\n", err)
		} else {
			log.Printf("Unable to run migrations: %v\n", err)
		}
	}
}
