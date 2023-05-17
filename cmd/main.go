package main

import (
	"fmt"

	"github.com/andremandrade/raiway/migrator"
)

const (
	profile        = ""
	database       = ""
	engine         = ""
	migrationsPath = "example/migrations"
)

func main() {

	migrationStatus_, setupErr := migrator.Setup(profile, database, engine, migrationsPath)

	if setupErr != nil {
		panic(setupErr)
	}

	migrationStatus := *migrationStatus_
	fmt.Println("\n=== Migration Status ===")
	fmt.Println(" - Current migration: ", migrationStatus.CurrentMigrationId)
	fmt.Println(" - Needed migrations:")
	for id := range migrationStatus.NextMigrations {
		fmt.Println("   - ", id)
	}

}
