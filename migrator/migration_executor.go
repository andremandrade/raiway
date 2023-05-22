package migrator

import (
	"fmt"

	"github.com/andremandrade/raiway/database"
)

func Execute(operation MigrationOperation) error {

	switch operation.Type {
	case LoadCsv:
		opExecError := executeLoadCSV(operation)
		if opExecError != nil {
			return fmt.Errorf(":Execute%w", opExecError)
		}
		return nil
	case LoadModels:
		opExecError := executeLoadModels(operation)
		if opExecError != nil {
			return fmt.Errorf(":Execute%w", opExecError)
		}
		return nil
	case DeleteModels:
		opExecError := executeDeleteModels(operation)
		if opExecError != nil {
			return fmt.Errorf(":Execute%w", opExecError)
		}
		return nil
	case Update:
		opExecError := executeUpdate(operation)
		if opExecError != nil {
			return fmt.Errorf(":Execute%w", opExecError)
		}
		return nil
	default:
		return fmt.Errorf(":Execute: Operation is not suported or is not implemented yet: %v", operation)
	}
}

func executeLoadCSV(operation MigrationOperation) error {
	validationError := validateLoadCSV(operation)
	if validationError != nil {
		return fmt.Errorf(":executeLoadCSV%w", validationError)
	}

	loadCsvError := database.LoadCsv(operation.ModelName, operation.FilePath,
		operation.Delimiter, operation.Quotechar, operation.Escapechar, operation.Scheme)
	if loadCsvError != nil {
		return fmt.Errorf(":executeLoadCSV%w", loadCsvError)
	}
	return nil
}

func executeLoadModels(operation MigrationOperation) error {
	validationError := validateLoadModels(operation)
	if validationError != nil {
		return fmt.Errorf(":executeLoadModels%w", validationError)
	}

	loadModelsError := database.LoadModels(operation.Prefix, operation.Files)
	if loadModelsError != nil {
		return fmt.Errorf(":executeLoadModels%w", loadModelsError)
	}
	return nil
}

func executeDeleteModels(operation MigrationOperation) error {
	validationError := validateDeleteModels(operation)
	if validationError != nil {
		return fmt.Errorf(":executeDeleteModels%w", validationError)
	}

	deleteModelsError := database.DeleteModels(operation.Models)
	if deleteModelsError != nil {
		return fmt.Errorf(":executeDeleteModels%w", deleteModelsError)
	}
	return nil
}

func executeUpdate(operation MigrationOperation) error {
	validationError := validateUpdate(operation)
	if validationError != nil {
		return fmt.Errorf(":executeUpdate%w", validationError)
	}

	if operation.Query != "" {
		_, updateError := database.Query(operation.Query, false)
		if updateError != nil {
			return fmt.Errorf(":executeUpdateInline%w", updateError)
		}
	} else if operation.FilePath != "" {
		_, updateError := database.QueryFromFile(operation.FilePath, false)
		if updateError != nil {
			return fmt.Errorf(":executeUpdateFromFile%w", updateError)
		}
	}
	return nil
}
