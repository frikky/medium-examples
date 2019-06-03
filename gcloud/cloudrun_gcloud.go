package main

import (
	"context"
	cloudrun "google.golang.org/api/run/v1"
	"log"
)

// https://cloud.google.com/run/docs/reference/rest/
func main() {
	ctx := context.Background()
	// Create a client like anywhere else
	service, err := cloudrun.NewService(ctx)
	if err != nil {
		log.Fatalf("Error creating cloudrun client: %s", err)
	}
	log.Printf("%#v", service)

	// POST https://run.googleapis.com/v1alpha1/{parent}/services

	newservice := cloudrun.NewProjectsService(service)
	log.Printf("%#v", newservice)

	// Time to go through struct hell again :)
	projectsLocationsService := cloudrun.NewProjectsLocationsService(service)

	// Wtf do I even do here
	projectsLocationsGetCall := projectsLocationsService.Get("projects/shuffle-241517/locations/us-central1/services")
	log.Printf("%#v", projectsLocationsGetCall)

	// googleapi.CallOption
	resp, err := projectsLocationsGetCall.Do()
	if err != nil {
		log.Println(err)
	}

	//log.Printf("%#v", projectsLocationsGetCall)
	//location, err := projectsLocationsGetCall.Do()
	//log.Println("%#v", location)
	//if err != nil {
	//	log.Fatalf("Locationgetcall issue: %s", err)
	//}
}
