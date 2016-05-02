package main

import (
  "github.com/nanobox-io/golang-scribble"
)

// Function to write to the database
func writeToDatabase(location string, elementID string, element interface{}) {
	// Open Database File
	db, err := scribble.New("./data", nil)
	if err != nil {
		panic(err)
	}
	// Write to the Database
	if err := db.Write(location, elementID, element); err != nil {
		panic(err)
	}
}

// Function to read from database
func readFromDatabase(location string, elementID string) (interface{}, error) {
	var element interface{}
	db, err := scribble.New("./data", nil)
	if err != nil {
		panic(err)
	}
	if err := db.Read(location, elementID, &element); err != nil {
		return element, err
	}
	return element, nil
}
