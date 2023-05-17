package migrator

import (
	"fmt"

	raiDatabase "github.com/andremandrade/raiway/database"
)

var migrationScripts map[int]MigrationScript

func Setup(profile, database, engine, migrationsPath string) (*MigrationStatus, error) {

	dbConnectionErr := raiDatabase.Connect(profile)

	if dbConnectionErr != nil {
		return nil, dbConnectionErr
	}

	raiDatabase.SetDefaultDatabaseAndEngine(database, engine)

	migrationScripts, err := LoadLocalMigrations(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}

	migrationStatus, err := GetMigrationStatus(migrationScripts)

	if err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}

	return migrationStatus, nil
}

func Migrate() error {
	return nil
}
