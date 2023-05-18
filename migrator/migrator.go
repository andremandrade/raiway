package migrator

import (
	"fmt"

	raiDatabase "github.com/andremandrade/raiway/database"
)

var allMigrationScripts map[int]MigrationScript
var currentMigrationStatus *MigrationStatus

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
	allMigrationScripts = migrationScripts
	migrationStatus, err := GetMigrationStatus(migrationScripts)

	if err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}
	currentMigrationStatus = migrationStatus
	return migrationStatus, nil
}

func GetLocalMigrationScripts() map[int]MigrationScript {
	return allMigrationScripts
}

func Migrate() error {
	fmt.Println("= Starting migration...")
	for mig_id, migrationScript := range *currentMigrationStatus.NextMigrations {
		fmt.Println("  * Migration script #", mig_id, " initiated...")
		for opId, op := range migrationScript {
			fmt.Println("     * Operation #", opId, " - ", op.Type, op.Name)
			execErr := Execute(op)
			if execErr != nil {
				fmt.Print("        ")
				fmt.Println(":Migrate: Operation execution failed: ", execErr)
			} else {
				fmt.Print("        ")
				fmt.Println(":Migrate: Operation execution succeeded")
			}
		}
	}
	return nil
}
