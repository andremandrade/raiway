package migrator

import (
	"fmt"

	raiDatabase "github.com/andremandrade/raiway/database"
	"github.com/relationalai/rai-sdk-go/rai"
)

type MigrationStatus struct {
	CurrentMigrationId int
	NextMigrations     map[int]MigrationScript
}

func GetMigrationStatus(migrationScripts map[int]MigrationScript) (*MigrationStatus, error) {
	tranx, err := raiDatabase.Query(CurrentMigrationQuery)
	if err != nil {
		return nil, fmt.Errorf("migrator.getMigrationStatus: %w", err)
	}

	noMigrations := &MigrationStatus{
		CurrentMigrationId: -1,
		NextMigrations:     migrationScripts,
	}
	if len(tranx.Problems) > 0 {
		problems := tranx.Problems
		for _, iProb := range problems {
			problem, typeInferenceSucceed := iProb.(rai.ClientProblem)
			if !typeInferenceSucceed {
				return nil, fmt.Errorf("migrator.getMigrationStatus: type inference failed for rai.ClientProblem")
			}
			if typeInferenceSucceed && problem.ErrorCode == "UNDEFINED" {
				return noMigrations, nil
			} else {
				return nil, fmt.Errorf("migrator.getMigrationStatus: unknown error - %v", problem)
			}
		}
	}

	if len(tranx.Results) == 0 {
		return noMigrations, nil
	}

	lastMigrationID := int(tranx.Results[0].Table[0].(float64))
	nextMigrations := make(map[int]MigrationScript)

	for id, migrationScript := range migrationScripts {
		if id > lastMigrationID {
			nextMigrations[id] = migrationScript
		}
	}

	migrationStatus := MigrationStatus{
		CurrentMigrationId: lastMigrationID,
		NextMigrations:     nextMigrations,
	}

	return &migrationStatus, nil
}

const CurrentMigrationQuery = `//beginrel
	def output = max[id: raiway_migrations(_, id, _)]
//endrel`
