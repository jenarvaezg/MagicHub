package db

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/zebresel-com/mongodm"
)

const (
	// DatabaseName holds the name of the database
	DatabaseName = "magichub"
)

func getLocalsMap() map[string]map[string]string {
	file, err := ioutil.ReadFile("locals.json")

	if err != nil {
		log.Panicf("File error: %v\n", err)
	}

	var localMap map[string]map[string]string
	json.Unmarshal(file, &localMap)

	return localMap
}

var conn *mongodm.Connection

// GetMongoConnection returns a mongodm connection to a mongodb database
func GetMongoConnection() *mongodm.Connection {
	if conn != nil {
		return conn
	}
	var err error

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		mongoURL = "127.0.0.1"
	}

	dbConfig := &mongodm.Config{
		DatabaseHosts: []string{mongoURL},
		DatabaseName:  DatabaseName,
		Locals:        getLocalsMap()["en-US"],
	}

	connection, err := mongodm.Connect(dbConfig)

	if err != nil {
		log.Panicf("Database connection error: %v", err)
	}

	conn = connection

	return conn
}
