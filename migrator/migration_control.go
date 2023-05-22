package migrator

import (
	"fmt"

	raiDatabase "github.com/andremandrade/raiway/database"
)

type MigrationStatus struct {
	DatabaseVersion      int
	InitFile             InitFile
	RecommendedMigration *MigrationFile
}

func GetMigrationStatus(initFile InitFile, migrationFiles []MigrationFile) (*MigrationStatus, error) {

	tranx, err := raiDatabase.Query(GetDatabaseVersionQuery, false)
	if err != nil {
		return nil, fmt.Errorf(":GetMigrationStatus: %w", err)
	}

	if len(tranx.Problems) > 0 {
		return nil, fmt.Errorf(`:GetMigrationStatus:
	RAI transaction Results: %v
	RAI transaction Problems: %v
	- Expected raiway:db_version relation not found
	- Run "raiway --init" to initialize a versioned database`, tranx.Results, tranx.Problems)
	}

	if len(tranx.Results) == 0 {
		return nil, fmt.Errorf(`:GetMigrationStatus:
	- Expected raiway:db_version relation not found
	- Run "raiway --init" to initialize a versioned database`)
	}

	databaseVersion := int(tranx.Results[0].Table[0].(float64))

	migrationStatus := MigrationStatus{
		DatabaseVersion: databaseVersion,
		InitFile:        initFile,
	}

	for fileIndex := range migrationFiles {
		migrationFile := migrationFiles[fileIndex]
		if migrationFile.SourceVersion == databaseVersion &&
			migrationFile.TargetVersion == initFile.Version {

			migrationStatus.RecommendedMigration = &migrationFile
		}
	}
	return &migrationStatus, nil
}

const GetDatabaseVersionQuery = `//beginrel
	def output = raiway:db_version
//endrel`
