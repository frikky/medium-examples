package main

import (
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/firestore"
)

var client *firestore.Client

type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"updateTime"`
	Fields     User      `json:"fields"`
}

// This is our self-defined fields.
// FirestoreEvent.Value.Fields = User
type User struct {
	Username   string `json:"userId"`
	Email      string `json:"email"`
	DateEdited int64  `json:"date_edited"`
}

// Simple init to have a firestore client available
func init() {
	ctx := context.Background()
	var err error
	client, err = firestore.NewClient(ctx, "medium-77273")
	if err != nil {
		log.Fatalf("Firestore: %v", err)
	}
}

// Handles the rollback to a previous document
func handleRollback(ctx context.Context, e FirestoreEvent) error {
	return errors.New("Should have rolled back to a previous version")
}

// The function that runs with the cloud function itself
func HandleUserChange(ctx context.Context, e FirestoreEvent) error {
	log.Printf("%#v", e)

	newFields := e.Value.Fields
	oldFields := e.OldValue.Fields
	// As our goal is simply to check if the username has changed
	if newFields.Username == oldFields.Username {
		return handleRollback(ctx, e)
	}

	// Check if the email is the same as previously
	if newFields.Email != oldFields.Email {
		return handleRollback(ctx, e)
	}

	// Check if the timestamp is older than old write to db
	// This is an int64 for storage reasons
	if newFields.DateEdited <= oldFields.DateEdited {
		return handleRollback(ctx, e)
	}

	// Successful
	return nil
}

func main() {
	ctx := context.Background()
	firestoreEvent := FirestoreEvent{}
	err := HandleUserChange(ctx, firestoreEvent)
	if err != nil {
		log.Printf("Err: %s", err)
	}
}
