package migrator

import (
	"fmt"

	"github.com/andremandrade/raiway/database"
)

func Execute(operation MigrationOperation) error {

	switch operation.Type {
	case LoadCsv:
		validationError := validateLoadCSV(operation)
		if validationError != nil {
			return fmt.Errorf(":Execute%w", validationError)
		}

		opExecError := executeLoadCSV(operation)
		if opExecError != nil {
			return fmt.Errorf(":Execute%w", opExecError)
		}
		return nil
	default:
		return fmt.Errorf(":Execute: Operation is not suported or is not implemented yet: %v", operation)
	}
}

func executeLoadCSV(operation MigrationOperation) error {
	loadCsvError := database.LoadCsv(operation.ModelName, operation.FilePath,
		operation.Delimiter, operation.Quotechar, operation.Escapechar, operation.Scheme)
	if loadCsvError != nil {
		return fmt.Errorf(":executeLoadCSV%w", loadCsvError)
	}
	return nil
}
