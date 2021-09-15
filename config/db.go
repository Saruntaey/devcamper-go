package config

import (
	"fmt"
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
		fmt.Printf("Database connection error: %v", err)
	}
	fmt.Printf("Connected to mongodb at %s\n", uri)
	return connection
}
