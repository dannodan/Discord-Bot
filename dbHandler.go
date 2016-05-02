package main

import (
  // "fmt"
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
func readFromDatabase(location string, elementID string) (map[string]interface{}, error) {
	var element map[string]interface{}
	db, err := scribble.New("./data", nil)
	if err != nil {
		panic(err)
	}
	if err := db.Read(location, elementID, &element); err != nil {
		return element, err
	}
	return element, nil
}

// Function to update in the Database
func updateToDatabase(location string, elementID string, arguments map[string]string) {
  var element map[string]interface{}
  // Open Database File
  db, err := scribble.New("./data", nil)
  if err != nil {
    panic(err)
  }
  // Get the Element
  if err := db.Read(location, elementID, &element); err != nil {
    panic(err)
  }
  // Change the attributes in Element
  for key, value := range arguments {
    element[key] = value
  }
  // Write to the Database
  if err := db.Write(location, elementID, element); err != nil {
    panic(err)
  }
}
