package infra

import (
	"eventstore-intro/eventstore/config"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"os"
)

func ConnectToEventstoreDB() (*esdb.Client, error) {
	connectionString := os.Getenv(config.EVENT_STORE_CONNECTION_STRING)
	settings, err := esdb.ParseConnectionString(connectionString)
	if err != nil {
		panic(err)
	}
	db, _ := esdb.NewClient(settings)
	return db, nil
}
