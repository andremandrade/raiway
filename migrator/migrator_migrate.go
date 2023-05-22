package migrator

import (
	"fmt"

	"github.com/andremandrade/raiway/database"
)

const insertMigrationRecordQuery = `//beginrel
def insert:raiway = :migration, %d, %d, datetime_now
//endrel`

func ExecuteMigration(initFile InitFile, migrationFiles []MigrationFile) error {

	migrationStatus, migrationStatusError := CheckMigrationStatus()

	if migrationStatusError != nil {
		return fmt.Errorf(":ExecuteMigration%w", migrationStatusError)
	}

	if migrationStatus.DatabaseVersion >= initFile.Version {
		return fmt.Errorf(":ExecuteMigration: Can not migrate! The database version is higher or equal to the init file version")
	}

	if migrationStatus.RecommendedMigration == nil {
		return fmt.Errorf(":ExecuteMigration: Can not migrate! No recommended migration was found (databaseVersion should be equals to targetVersion)")
	}
	fmt.Println("  Database is going to be migrated from",
		migrationStatus.RecommendedMigration.SourceVersion, "to",
		migrationStatus.RecommendedMigration.TargetVersion)
	for id, operation := range migrationStatus.RecommendedMigration.Script {
		execError := Execute(operation)
		if execError != nil {
			return fmt.Errorf(":ExecuteMigration%w", execError)
		}
		fmt.Println("  ✓ ", id, "- Successfully executed ", operation.Type, operation.Name)
	}
	insertMigrationRecord(migrationStatus.RecommendedMigration.SourceVersion, migrationStatus.RecommendedMigration.TargetVersion)
	fmt.Println("  ✓ Migration record created")
	updateDatabaseVersion(migrationStatus.RecommendedMigration.TargetVersion)
	fmt.Println("  ✓ Database version upgraded to version", migrationStatus.RecommendedMigration.TargetVersion)
	return nil
}

func insertMigrationRecord(sourceVersion, targetVersion int) error {
	query := fmt.Sprintf(insertMigrationRecordQuery, sourceVersion, targetVersion)
	_, queryExecError := database.Query(query, false)
	if queryExecError != nil {
		return fmt.Errorf(":insertMigrationRecord%w", queryExecError)
	}
	return nil
}
