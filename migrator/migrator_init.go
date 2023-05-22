package migrator

import (
	"fmt"

	"github.com/andremandrade/raiway/database"
)

const deleteExistingVersionQuery = `//beginrel
def delete:raiway = :db_version, raiway:db_version
//endrel`
const insertVersionQuery = `//beginrel
def insert:raiway = :db_version, %d
//endrel`

func ExecuteInit(initFile InitFile) error {
	fmt.Println("  Database is going to be initilized with version", initFile.Version)
	for id, operation := range initFile.Script {
		execError := Execute(operation)
		if execError != nil {
			return fmt.Errorf(":ExecInit%w", execError)
		}
		fmt.Println("  ✓ ", id, "- Successfully executed ", operation.Type, operation.Name)
	}
	updateDatabaseVersion(initFile.Version)
	return nil
}

func updateDatabaseVersion(version int) error {
	baseRelations, dbError := database.GetBaseRelations()
	if dbError != nil {
		return fmt.Errorf(":updateDatabaseVersion%w", dbError)
	}
	raiwayEDBexists := false
	for _, baseRelation := range baseRelations {
		if baseRelation == "raiway" {
			raiwayEDBexists = true
		}
	}
	updateVersionQuery := ""
	if !raiwayEDBexists {
		updateVersionQuery += deleteExistingVersionQuery
	}
	updateVersionQuery += fmt.Sprintf(insertVersionQuery, version)
	_, queryExecError := database.Query(updateVersionQuery, false)
	if queryExecError != nil {
		return fmt.Errorf(":updateDatabaseVersion%w", queryExecError)
	}
	return nil
}
