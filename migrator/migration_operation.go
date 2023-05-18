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
		return fmt.Errorf(":validateLoadCSV: Invalid operation type")
	}
	if op.HostMode != Local && op.HostMode != Cloud {
		return fmt.Errorf(":validateLoadCSV: Invalid hostMode")
	}
	if op.FilePath == "" || op.ModelName == "" {
		return fmt.Errorf(":validateLoadCSV: Missing required property: filePath or modelName")
	}
	return nil
}
