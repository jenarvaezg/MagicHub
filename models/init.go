package models

import (
	"log"

	"github.com/go-bongo/bongo"
	"github.com/jenarvaezg/magicbox/utils"
)

var connection *bongo.Connection

//Model interface is an interface for CRUD objects
type Model interface {
	Save() error
	Delete() error
	Update(updateMap utils.JSONMap) error
}

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
