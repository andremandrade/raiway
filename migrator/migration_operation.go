package migrator

import "fmt"

type MigrationOperation struct {
	//all
	Type OperationType `yaml:"type"`
	//load-csv
	HostMode HostMode `yaml:"hostMode"`
	// load-csv load-models update
	ModelName string `yaml:"modelName,omitempty"`
	// load-csv load-models update
	FilePath string `yaml:"filePath,omitempty"`
	// load-models
	Delimiter  string            `yaml:"delimiter,omitempty"`
	Quotechar  string            `yaml:"quotechar,omitempty"`
	Escapechar string            `yaml:"escapechar,omitempty"`
	Scheme     map[string]string `yaml:"scheme,omitempty"`
	//all (optional)
	Name string `yaml:"name,omitempty"`
	// load-models
	Files []string `yaml:"files,omitempty"`
	// load-models
	Prefix string `yaml:"prefix,omitempty"`
	// load-models
	Models []string `yaml:"models,omitempty"`
	// update
	Query string `yaml:"query,omitempty"`
}

type OperationType string

const (
	LoadCsv      OperationType = "load-csv"
	LoadModels   OperationType = "load-models"
	DeleteModels OperationType = "delete-models"
	EnableICs    OperationType = "enable-ics"
	DisableICs   OperationType = "disable-ics"
	Update       OperationType = "update"
)

type HostMode string

const (
	Local HostMode = "local"
	Cloud HostMode = "cloud"
)

func validateLoadCSV(op MigrationOperation) error {
	if op.Type != LoadCsv {
		return fmt.Errorf(":validateLoadCSV: Operation type should be %s", LoadCsv)
	}
	if op.HostMode != Local && op.HostMode != Cloud {
		return fmt.Errorf(":validateLoadCSV: Invalid hostMode")
	}
	if op.FilePath == "" || op.ModelName == "" {
		return fmt.Errorf(":validateLoadCSV: Missing required property: filePath or modelName")
	}
	return nil
}

func validateLoadModels(op MigrationOperation) error {
	if op.Type != LoadModels {
		return fmt.Errorf(":validateLoadModels: Operation type should be %s", LoadModels)
	}
	if op.Files == nil || len(op.Files) == 0 {
		return fmt.Errorf(":validateLoadModels: Missing required property: files")
	}

	for _, f := range op.Files {
		if f == "" {
			return fmt.Errorf(":validateLoadModels: files element can not be empty")
		}
	}

	if op.Prefix == "" {
		return fmt.Errorf(":validateLoadModels: Missing required property: prefix")
	}
	return nil
}

func validateDeleteModels(op MigrationOperation) error {
	if op.Type != DeleteModels {
		return fmt.Errorf(":validateDeleteModels: Operation type should be %s", DeleteModels)
	}
	if op.Models == nil || len(op.Models) == 0 {
		return fmt.Errorf(":validateDeleteModels: Missing required property: models")
	}

	for _, m := range op.Models {
		if m == "" {
			return fmt.Errorf(":validateDeleteModels: models element can not be empty")
		}
	}

	return nil
}

func validateUpdate(op MigrationOperation) error {
	if op.Type != Update {
		return fmt.Errorf(":validateUpdate: Operation type should be %s", Update)
	}
	if op.FilePath == "" && op.Query == "" {
		return fmt.Errorf(":validateUpdate: Missing required property: filePath or query")
	}
	return nil
}
