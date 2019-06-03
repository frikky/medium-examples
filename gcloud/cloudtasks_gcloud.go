package main

import (
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
	"log"
)

func createTask(ctx context.Context, client *cloudtasks.Client, parent string) {
	// Define some endpoint you want the data to hit from
	url := "/api/test"

	// Nested structs. Just mapped them like this so it's actually readable
	var appEngineHttpRequest *taskspb.AppEngineHttpRequest = &taskspb.AppEngineHttpRequest{
		HttpMethod:  taskspb.HttpMethod_GET,
		RelativeUri: url,
	}

	var appeng *taskspb.Task_AppEngineHttpRequest = &taskspb.Task_AppEngineHttpRequest{
		AppEngineHttpRequest: appEngineHttpRequest,
	}

	var task *taskspb.Task = &taskspb.Task{
		PayloadType: appeng,
	}

	// Structs added into the last struct which creates the task
	req := &taskspb.CreateTaskRequest{
		Parent: parent,
		Task:   task,
	}

	ret, err := client.CreateTask(ctx, req)
	if err != nil {
		log.Printf("Error creating task: %s", err)
		return
	}

	log.Printf("%#v", ret)

}

func listAllTasks(ctx context.Context, client *cloudtasks.Client, parent string) {
	// Makes a struct to map
	req := &taskspb.ListTasksRequest{
		Parent: parent,
	}

	// Returns an iterator over the parent tasks and counts
	ret := client.ListTasks(ctx, req)
	cnt := 0
	for {
		_, err := ret.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Printf("Error in iterator: %s", err)
			break
		}

		cnt += 1
	}

	log.Printf("Current amount of tasks: %d", cnt)
}

func main() {
	// Define the client
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Fatalf("Error creating cloudtask client: %s", err)
	}

	// Set the projectId, location and queuename for the specific request
	projectID := "yourprojectnamehere"
	location := "europe-west3"
	queuename := "myqueue"
	var formattedParent string = fmt.Sprintf("projects/%s/locations/%s/queues/myqueue", projectID, location, queuename)

	// Creates a task
	createTask(ctx, client, formattedParent)
	listAllTasks(ctx, client, formattedParent)
}
