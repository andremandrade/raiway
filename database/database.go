package database

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

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

func LoadCsv(relation, filePath, delimiter, quotechar, escapechar string, csvSchema map[string]string) error {

	file, fileReadError := ioutil.ReadFile(filePath)
	if fileReadError != nil {
		return fmt.Errorf(":LoadCsv:fileReadError:%w", fileReadError)
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
		return fmt.Errorf(":LoadCsv:raiClient.LoadCSV:%w", loadCsvError)
	}

	tranxError := checkTransactionSuccess(tranx)
	if tranxError != nil {
		return fmt.Errorf(":LoadCsv%w", tranxError)
	}
	return nil
}

func LoadModels(prefix string, files []string) error {
	filesIoReaders := make(map[string]io.Reader)
	for _, filePath := range files {
		file, fileReadError := ioutil.ReadFile(filePath)
		if fileReadError != nil {
			return fmt.Errorf(":LoadModels:fileReadError:%w", fileReadError)
		}

		ioReader := bytes.NewReader(file)
		splittedFilePath := strings.Split(filePath, "/")
		modelName := strings.TrimSuffix(fmt.Sprint(prefix, "/", splittedFilePath[len(splittedFilePath)-1]), ".rel")
		filesIoReaders[modelName] = ioReader
	}

	tranx, loadModelsError := raiClient.LoadModels(database, engine, filesIoReaders)

	if loadModelsError != nil {
		return fmt.Errorf(":LoadModels:raiClient.LoadModels:%w", loadModelsError)
	}

	tranxError := checkTransactionSuccess(tranx)
	if tranxError != nil {
		return fmt.Errorf(":LoadModels%w", tranxError)
	}
	return nil
}

func DeleteModels(models []string) error {
	tranx, loadModelsError := raiClient.DeleteModels(database, engine, models)

	if loadModelsError != nil {
		return fmt.Errorf(":DeleteModels:raiClient.LoadModels:%w", loadModelsError)
	}

	tranxError := checkTransactionSuccess(tranx)
	if tranxError != nil {
		return fmt.Errorf(":DeleteModels%w", tranxError)
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
