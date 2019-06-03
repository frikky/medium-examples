package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"log"
)

// Create an item we're gonna put in and remove
// Uses datastore and json tags to map it directly
type Data struct {
	Key string `datastore:"key" json:"key"`
}

// Puts some data from the struct Data into the database
func putDatastore(ctx context.Context, client *datastore.Client, dbname string, data Data) error {
	// Make a key to map to datastore
	datastoreKey := datastore.NameKey(dbname, data.Key, nil)

	// Adds the key described above wit hthe data from datastoreKey
	if _, err := client.Put(ctx, datastoreKey, &data); err != nil {
		log.Fatalf("Error adding testdata to %s: %s", dbname, err)
		return err
	}

	return nil
}

// Gets the data back from the datastore
func getDatastore(ctx context.Context, client *datastore.Client, dbname string, identifier string) (*Data, error) {
	// Defines the key
	datastoreKey := datastore.NameKey(dbname, identifier, nil)

	// Creates an empty variable of struct Data, which we map the data back to
	newdata := &Data{}
	if err := client.Get(ctx, datastoreKey, newdata); err != nil {
		return &Data{}, err
	}

	return newdata, nil
}

func main() {
	// Describe the project
	projectID := "yourprojectnamehere"
	dbname := "medium-test"

	ctx := context.Background()

	// Create a client
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed setting up client")
	}

	// Create som json data to map to struct
	jsondata := `{
		"key": "qwertyuiopasdfghjkl"
	}`

	// Map the jsondata to the struct Data
	var structData Data
	if err := json.Unmarshal([]byte(jsondata), &structData); err != nil {
		log.Fatalf("Failed unmarshalling: %s", err)
	}

	// Puts the data described above in the datastore
	if err := putDatastore(ctx, client, dbname, structData); err != nil {
		log.Fatalf("Failed putting in datastore: %s", err)
	}

	// Gets the same data back from the datastore
	returnData, err := getDatastore(ctx, client, dbname, structData.Key)
	if err != nil {
		log.Fatalf("Failed getting from datastore: %s", err)
	}

	// Print with some extra value
	log.Printf("%#v", returnData)
}
