package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

func ConnDB() *mgo.Session {
	uri := os.Getenv("MONGO_URI")
	s, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Connected to mongodb at %s\n", uri)
	return s
}
