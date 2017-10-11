package models

import (
	"log"

	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

var connection *bongo.Connection

var emptyUser = bson.ObjectId(0)

func connectToMongo() *bongo.Connection {
	config := &bongo.Config{
		ConnectionString: "localhost",
		Database:         "bongotest",
	}
	var err error
	connection, err = bongo.Connect(config)

	if err != nil {
		log.Fatal(err)
	}
	return connection

}

func setupCollections() {
	boxCollection = connection.Collection("box")
	userCollection = connection.Collection("user")
}

func init() {
	connectToMongo()
	setupCollections()
}
