package migrator

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type MigrationOperation string

const (
	LoadCsv      MigrationOperation = "load-csv"
	LoadModels   MigrationOperation = "load-models"
	DeleteModels MigrationOperation = "delete-models"
	EnableICs    MigrationOperation = "enable-ics"
	DisableICs   MigrationOperation = "disable-ics"
	Update       MigrationOperation = "update"
)

type MigrationScript []struct {
	//all
	Type MigrationOperation `yaml:"type"`
	// load-csv load-models update
	FilePath string `yaml:"filePath,omitempty"`
	// load-models
	Scheme map[string]string `yaml:"scheme,omitempty"`
	//all (optional)
	Name string `yaml:"name,omitempty"`
	// load-models
	Files []string `yaml:"files,omitempty"`
	// load-models
	Prefix string `yaml:"prefix,omitempty"`
	// update
	Query string `yaml:"query,omitempty"`
}

func LoadLocalMigrations(migrationsPath string) (map[int]MigrationScript, error) {

	migrationFiles, migrationFilesError := getMigrationFiles(migrationsPath)

	if migrationFilesError != nil {
		return nil, migrationFilesError
	}

	migrationScripts := make(map[int]MigrationScript)

	fmt.Println("Found migration files")
	fmt.Println("ID | Name ")
	for id, migrationFile := range migrationFiles {
		migrationScript, createMigrationErr := createMigration(migrationsPath, migrationFile)
		if createMigrationErr != nil {
			return nil, fmt.Errorf("load local migrations:  %w", createMigrationErr)
		}
		migrationScripts[id] = migrationScript
		fmt.Println(id, "|", migrationFile)
	}

	return migrationScripts, nil
}

func createMigration(migrationsPath, migrationFile string) (MigrationScript, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", migrationsPath, migrationFile))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var migrationScript MigrationScript

	err = yaml.Unmarshal(data, &migrationScript)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling migration file %s :  %w", migrationFile, err)
	}

	return migrationScript, nil
}

func getMigrationFiles(migrationsPath string) (map[int]string, error) {
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir %s :  %w", migrationsPath, err)
	}

	migrations := make(map[int]string)

	for _, file := range files {
		if !file.IsDir() {
			if !strings.HasSuffix(file.Name(), DotYAML) {
				fmt.Printf("WARNING - %s was ignored because '.yaml' extension was not found", file.Name())
			}
			migrationNumber, convertErr := strconv.Atoi(strings.TrimSuffix(file.Name(), DotYAML))

			if convertErr != nil {
				return nil, fmt.Errorf("migration file name %s is not a number:  %w", file.Name(), err)
			}
			migrations[migrationNumber] = file.Name()
		}
	}

	sortedMigrationNumbers := make([]int, 0, len(migrations))
	for id := range migrations {
		sortedMigrationNumbers = append(sortedMigrationNumbers, id)
	}

	sort.Ints(sortedMigrationNumbers)

	sortedMigrations := make(map[int]string)
	for _, id := range sortedMigrationNumbers {
		sortedMigrations[id] = migrations[id]
	}

	return sortedMigrations, nil
}

const DotYAML = ".yaml"
