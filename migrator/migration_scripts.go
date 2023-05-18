package migrator

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type MigrationScript []MigrationOperation

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
			return nil, fmt.Errorf(":Load local migrations%w", createMigrationErr)
		}
		migrationScripts[id] = migrationScript
		fmt.Println(id, "|", migrationFile)
	}

	return migrationScripts, nil
}

func createMigration(migrationsPath, migrationFile string) (MigrationScript, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", migrationsPath, migrationFile))
	if err != nil {
		return nil, fmt.Errorf(":createMigration:readFile:%w", err)
	}

	var migrationScript MigrationScript

	err = yaml.Unmarshal(data, &migrationScript)
	if err != nil {
		return nil, fmt.Errorf(":createMigration:unmarshalling migration file %s:%w", migrationFile, err)
	}

	return migrationScript, nil
}

func getMigrationFiles(migrationsPath string) (map[int]string, error) {
	files, fileReadError := ioutil.ReadDir(migrationsPath)
	if fileReadError != nil {
		return nil, fmt.Errorf(":Read migrations dir %s%w", migrationsPath, fileReadError)
	}

	migrations := make(map[int]string)

	for _, file := range files {
		if !file.IsDir() {
			if !strings.HasSuffix(file.Name(), DotYAML) {
				fmt.Printf("WARNING - %s was ignored because '.yaml' extension was not found", file.Name())
				continue
			}
			migrationNumber, convertionErr := strconv.Atoi(strings.TrimSuffix(file.Name(), DotYAML))

			if convertionErr != nil {
				return nil, fmt.Errorf(":Migration file name %s is not a number:  %w", file.Name(), fileReadError)
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
