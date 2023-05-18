package database

import (
	"bytes"
	"fmt"
	"io/ioutil"

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

func LoadCsv(relation, fileLocation, delimiter, quotechar, escapechar string, csvSchema map[string]string) error {

	file, fileReadError := ioutil.ReadFile(fileLocation)
	if fileReadError != nil {
		return fmt.Errorf(":fileReadError:%w", fileReadError)
	}

	ioReader := bytes.NewReader(file)

	csvOpts := rai.CSVOptions{}

	if delimiter != "" {
		csvOpts.WithDelim([]rune(delimiter)[0])
	}
	if quotechar != "" {
		csvOpts.WithQuoteChar([]rune(quotechar)[0])
	}
	if escapechar != "" {
		csvOpts.WithEscapeChar([]rune(escapechar)[0])
	}
	if csvSchema != nil {
		csvOpts.WithSchema(csvSchema)
	}

	tranx, loadCsvError := raiClient.LoadCSV(database, engine, relation, ioReader, &csvOpts)

	if loadCsvError != nil {
		return fmt.Errorf(":loadCsv:raiClient.LoadCSV")
	}

	tranxError := checkTransactionSuccess(tranx)
	if tranxError != nil {
		return fmt.Errorf(":loadCsv%w", tranxError)
	}
	return nil
}

func checkTransactionSuccess(tranx *rai.TransactionResult) error {
	if tranx.Aborted {
		return fmt.Errorf(":checkTransactionSucces: Transaction was aborted")
	}
	if len(tranx.Problems) > 0 {
		return fmt.Errorf(":checkTransactionSucces: has problems %v", tranx.Problems)
	}
	return nil
}
