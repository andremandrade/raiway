package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/andremandrade/raiway/migrator"
)

func main() {

	var profile string
	var database string
	var engine string
	var migrationsPath string

	flag.StringVar(&profile, "p", "default", "[--profile]")
	flag.StringVar(&profile, "profile", "default", "RAI profile name")

	flag.StringVar(&database, "d", "", "[--database]")
	flag.StringVar(&database, "database", "", "RAI database name")

	flag.StringVar(&engine, "e", "", "[--engine]")
	flag.StringVar(&engine, "engine", "", "RAI engine name")

	flag.StringVar(&migrationsPath, "f", "./migrations", "[--folder]")
	flag.StringVar(&migrationsPath, "folder", "./migrations", "Migrations folder path")

	var init bool
	var migrate bool

	flag.BoolVar(&init, "i", false, "[--init]")
	flag.BoolVar(&init, "init", false, "Run database initiatization script")

	flag.BoolVar(&migrate, "m", false, "[--migrate]")
	flag.BoolVar(&migrate, "migrate", false, "Search and apply a migration if it exists")

	flag.Parse()
	// fmt.Println(init, migrate, profile, database, engine, migrationsPath)

	initFile, migrationFiles, setupErr := migrator.Setup(profile, database, engine, migrationsPath)

	if setupErr != nil {
		panic(fmt.Errorf(":raiway-cli%w", setupErr))
	}

	fmt.Println(`[RAi Database config]
    ✓ Succesfully connected`)

	printLocalConfig(initFile, migrationFiles)

	if init || migrate {
		if initFile == nil {
			fmt.Println(errors.New("[ERROR] User options can not be executed"))
			return
		}
		executeOptionsError := executeOptions(init, migrate, *initFile, migrationFiles)
		if executeOptionsError != nil {
			return
		}
	}
}

func executeOptions(init, migrate bool, initFile migrator.InitFile, migrationFiles []migrator.MigrationFile) error {
	if init {
		if migrate {
			fmt.Println("[warning] Argument --migrate is ignored because can't be used with --init")
		}
		execInitError := migrator.ExecuteInit(initFile)
		if execInitError != nil {
			return execInitError
		}
		fmt.Println("[Init] Database was successfully setup to the version", initFile.Version)
	}
	return nil
}

func printLocalConfig(initFile *migrator.InitFile, migrationFiles []migrator.MigrationFile) {
	fmt.Println("[RAiway local config]")

	if initFile != nil {
		fmt.Println("    ✓ init file found")
	} else {
		fmt.Println("    ! init file not found")
	}
	if len(migrationFiles) > 0 {
		fmt.Printf("    ✓ %d migration files found\n", len(migrationFiles))
	} else {
		fmt.Println("    ! No migration files")
	}
}
