package migrator

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

type InitFile struct {
	Version int             `yaml:"version"`
	Script  MigrationScript `yaml:"script"`
}
type MigrationFile struct {
	SourceVersion int             `yaml:"sourceVersion"`
	TargetVersion int             `yaml:"targetVersion"`
	Script        MigrationScript `yaml:"script"`
}

type MigrationScript []MigrationOperation

func LoadLocalMigrations(migrationsPath string) (*InitFile, []MigrationFile, error) {

	initFileName, migrationFilesName, migrationFilesError := getMigrationFiles(migrationsPath)

	if migrationFilesError != nil {
		return nil, nil, fmt.Errorf(":LoadLocalMigrations%w", migrationFilesError)
	}

	initFile := &InitFile{}
	if initFileName == "" {
		initFile = nil
	} else {
		initFileError := createInitFile(initFile, migrationsPath, initFileName)
		if initFileError != nil {
			return nil, nil, fmt.Errorf(":LoadLocalMigrations%w", initFileError)
		}
	}

	if len(migrationFilesName) == 0 {
		return initFile, nil, nil
	}
	migrationFiles := []MigrationFile{}
	for fileNameIndex := range migrationFilesName {
		migrationFile, migrationFileError := createMigrationFile(migrationsPath, migrationFilesName[fileNameIndex])
		if migrationFileError != nil {
			return initFile, migrationFiles, fmt.Errorf(":LoadLocalMigrations%w", migrationFileError)
		}
		migrationFiles = append(migrationFiles, *migrationFile)
	}

	return initFile, migrationFiles, nil
}

func createInitFile(initFile *InitFile, migrationsPath, fileName string) error {
	fileData, readFileError := ioutil.ReadFile(fmt.Sprintf("%s/%s", migrationsPath, fileName))
	if readFileError != nil {
		return fmt.Errorf(":createInitFile:readFile:%w", readFileError)
	}

	readFileError = yaml.Unmarshal(fileData, &initFile)
	if readFileError != nil {
		return fmt.Errorf(":createInitFile:unmarshalling init file %s:%w", fileName, readFileError)
	}

	return nil
}

func createMigrationFile(migrationsPath, fileName string) (*MigrationFile, error) {
	fileData, readFileError := ioutil.ReadFile(fmt.Sprintf("%s/%s", migrationsPath, fileName))
	if readFileError != nil {
		return nil, fmt.Errorf(":createMigrationFile:readFile:%w", readFileError)
	}

	var migrationFile MigrationFile
	unmarshalError := yaml.Unmarshal(fileData, &migrationFile)
	if unmarshalError != nil {
		return nil, fmt.Errorf(":createMigrationFile:unmarshalling migration file %s:%w", fileName, unmarshalError)
	}
	if migrationFile.SourceVersion >= migrationFile.TargetVersion {
		return nil, fmt.Errorf(":createMigrationFile: Error in %s : targetVersion should be greater than sourceVersion", fileName)
	}
	return &migrationFile, nil
}

func getMigrationFiles(migrationsPath string) (string, []string, error) {
	files, fileReadError := ioutil.ReadDir(migrationsPath)
	if fileReadError != nil {
		return "", nil, fmt.Errorf(": Read migrations dir %s : %w", migrationsPath, fileReadError)
	}

	var initFile string
	migrationFiles := []string{}

	for _, file := range files {
		if !file.IsDir() {
			if !strings.HasSuffix(file.Name(), DotYAML) && !strings.HasSuffix(file.Name(), DotYML) {
				fmt.Printf("WARNING - %s was ignored: invalid YAML extension", file.Name())
				continue
			}

			if file.Name() == init_yaml || file.Name() == init_yml {
				initFile = file.Name()
				continue
			}
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	return initFile, migrationFiles, nil
}

const DotYAML = ".yaml"
const DotYML = ".yml"
const init_yaml = "init.yaml"
const init_yml = "init.yml"
