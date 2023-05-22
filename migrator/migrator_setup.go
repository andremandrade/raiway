package migrator

import (
	"fmt"

	raiDatabase "github.com/andremandrade/raiway/database"
)

var isSetupExecuted = false
var allMigrations []MigrationFile
var currentInitFile *InitFile
var currentMigrationStatus *MigrationStatus

func Setup(profile, database, engine, migrationsPath string) (*InitFile, []MigrationFile, error) {

	dbConnectionErr := raiDatabase.Connect(profile)

	if dbConnectionErr != nil {
		return nil, nil, dbConnectionErr
	}

	raiDatabase.SetDefaultDatabaseAndEngine(database, engine)

	initFile, migrationFiles, localMigrationsError := LoadLocalMigrations(migrationsPath)
	if localMigrationsError != nil {
		return nil, nil, fmt.Errorf(":Setup%w", localMigrationsError)
	}

	currentInitFile = initFile
	allMigrations = migrationFiles
	isSetupExecuted = true
	return initFile, migrationFiles, nil
}

func CheckMigrationStatus() (*MigrationStatus, error) {
	if !isSetupExecuted {
		return nil, fmt.Errorf(":CheckMigrationStatus: migrator.Setup execution is required")
	}
	migrationStatus, migrationStatusError := GetMigrationStatus(*currentInitFile, allMigrations)

	if migrationStatusError != nil {
		return nil, fmt.Errorf("CheckMigrationStatus: %w", migrationStatusError)
	}
	currentMigrationStatus = migrationStatus
	return migrationStatus, nil
}

func GetLocalMigrationScripts() []MigrationFile {
	return allMigrations
}

func GetCurrentMigrationStatus() MigrationStatus {
	return *currentMigrationStatus
}
