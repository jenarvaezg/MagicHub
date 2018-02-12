package models

import (
	"log"
	"os"

	"github.com/go-bongo/bongo"
)

var connection *bongo.Connection

func connectToMongo() *bongo.Connection {

	mongoURL := os.Getenv("MONGO_URL")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	config := &bongo.Config{
		ConnectionString: mongoURL,
		Database:         mongoDatabase,
	}
	var err error
	log.Println("Connection to mongo")
	connection, err = bongo.Connect(config)
	log.Println("Connected to mongo")
	if err != nil {
		log.Fatal(err)
	}
	return connection

}

func setupCollections() {
	boxCollection = connection.Collection("box")
	userCollection = connection.Collection("user")
	log.Println("Collections ready")
}

func init() {
	connectToMongo()
	setupCollections()
}
