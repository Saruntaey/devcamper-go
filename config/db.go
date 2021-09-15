package config

import (
	"fmt"
	"log"
	"os"

	"github.com/zebresel-com/mongodm"
)

func ConnDB() *mongodm.Connection {
	uri := os.Getenv("MONGO_URI")
	dbConfig := &mongodm.Config{
		DatabaseHosts: []string{uri},
		DatabaseName:  os.Getenv("MONGO_DB"),
	}

	connection, err := mongodm.Connect(dbConfig)

	if err != nil {
		log.Fatalf("Database connection error: %v\n", err)
	}
	fmt.Printf("Connected to mongodb at %s\n", uri)
	return connection
}
