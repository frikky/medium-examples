package user

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/frikky/firestore-rollback-go"
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
	Username   rollback.StringValue  `json:"userId"`
	Email      rollback.StringValue  `json:"email"`
	DateEdited rollback.IntegerValue `json:"date_edited"`
}

// Simple init to have a firestore client available
func init() {
	ctx := context.Background()
	var err error
	// FIXME - add your username in here
	client, err = firestore.NewClient(ctx, "YOUR-PROJECT")
	if err != nil {
		log.Fatalf("Firestore: %v", err)
	}
}

// Handles the rollback of any previous document
func handleRollback(ctx context.Context, e FirestoreEvent) error {
	writtenData, firestoreReturn, err := rollback.Rollback(
		ctx,
		client,            // Your firestore client
		e.OldValue.Name,   // The path to roll back
		e.OldValue.Fields, // The parsed data you have
	)

	log.Printf("ROLLED BACK WITH DATA: %#v", writtenData)
	log.Printf("FIRESTORE RETURN: %#v", firestoreReturn)

	return err
}

// The function that runs with the cloud function itself
func HandleUserChange(ctx context.Context, e FirestoreEvent) error {
	// This is the data that's in the database itself
	newFields := e.Value.Fields
	oldFields := e.OldValue.Fields

	// Check if the email is the same as previously
	if newFields.Email.StringValue != oldFields.Email.StringValue {
		log.Printf("Bad email: %s - %s", newFields.Email.StringValue, oldFields.Email.StringValue)
		return handleRollback(ctx, e)
	}

	// As our goal is simply to check if the username has changed
	if newFields.Username.StringValue == oldFields.Username.StringValue {
		log.Printf("Bad username (same): %s - %s", newFields.Username.StringValue, oldFields.Username.StringValue)
		return handleRollback(ctx, e)
	}

	log.Printf("User successfully changed from %s to %s", newFields.Username.StringValue, oldFields.Username.StringValue)
	return nil
}
