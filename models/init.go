package models

import (
	"log"

	"github.com/go-bongo/bongo"
)

var connection *bongo.Connection

func init() {
	config := &bongo.Config{
		ConnectionString: "localhost",
		Database:         "bongotest",
	}
	var err error
	connection, err = bongo.Connect(config)

	if err != nil {
		log.Fatal(err)
	}
}
