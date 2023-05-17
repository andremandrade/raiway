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

	if len(tranx.Problems) > 0 {
		problems := tranx.Problems
		for _, iProb := range problems {
			problem, notErr := iProb.(rai.ClientProblem)
			if !notErr {
				return nil, fmt.Errorf("migrator.getMigrationStatus: type inference failed for rai.ClientProblem")
			}
			if notErr && problem.ErrorCode == "UNDEFINED" {
				return &MigrationStatus{
					CurrentMigrationId: -1,
					NextMigrations:     migrationScripts,
				}, nil
			} else {
				return nil, fmt.Errorf("migrator.getMigrationStatus: unknown error - %v", problem)
			}
		}
	}
	return nil, nil
}

const CurrentMigrationQuery = `//beginrel
	def output = raiway_migrations
//endrel`
