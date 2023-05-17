package database

import (
	"fmt"

	"github.com/relationalai/rai-sdk-go/rai"
)

var raiClient rai.Client
var database, engine string

func Connect(profile string) error {
	conn, connectionErr := rai.NewClientFromConfig(profile)

	if connectionErr != nil {
		return connectionErr
	}
	fmt.Println("RAI database connection stablished")
	raiClient = *conn
	return nil
}

func SetDefaultDatabaseAndEngine(defDatabase, defEngine string) {
	database = defDatabase
	engine = defEngine
}

func Query(query string) (*rai.TransactionAsyncResult, error) {
	result, err := raiClient.Execute(database, engine, query, nil, true)
	if err != nil {
		return nil, fmt.Errorf("database.query: %w", err)
	}
	return result, nil
}
