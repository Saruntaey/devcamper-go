package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/zebresel-com/mongodm"
)

func ConnDB() *mongodm.Connection {

	// Load prompt text when validate data before save fail
	file, err := ioutil.ReadFile("./config/locals.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var localMap map[string]map[string]string

	json.Unmarshal(file, &localMap)
	uri := os.Getenv("MONGO_URI")
	dbConfig := &mongodm.Config{
		DatabaseHosts: []string{uri},
		DatabaseName:  os.Getenv("MONGO_DB"),

		// Mount validation prompt text
		Locals: localMap["en-US"],
	}

	connection, err := mongodm.Connect(dbConfig)

	if err != nil {
		log.Fatalf("Database connection error: %v\n", err)
	}
	fmt.Printf("Connected to mongodb at %s\n", uri)
	return connection
}
