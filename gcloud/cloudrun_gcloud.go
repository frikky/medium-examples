package main

import (
	"context"
	"fmt"
	cloudrun "google.golang.org/api/run/v1alpha1"
	"log"
)

func getAllLocations(projectsLocationsService *cloudrun.ProjectsLocationsService) ([]string, error) {
	// List locations
	// Make a request, then do the request
	list := projectsLocationsService.List(fmt.Sprintf("projects/shuffle-241517"))
	ret, err := list.Do()

	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	locationNames := []string{}
	for _, item := range ret.Locations {
		locationNames = append(locationNames, item.Name)
	}

	return locationNames, nil

}

// https://cloud.google.com/run/docs/reference/rest/
func main() {
	ctx := context.Background()

	// Create a service client like anywhere else
	apiservice, err := cloudrun.NewService(ctx)
	if err != nil {
		log.Fatalf("Error creating cloudrun service client: %s", err)
	}

	// Gets all available locations
	projectsLocationsService := cloudrun.NewProjectsLocationsService(apiservice)
	allLocations, err := getAllLocations(projectsLocationsService)
	if err != nil {
		log.Fatalf("Error getting locations: %s", err)
	}

	// Define an image to use
	imagename := "gcr.io/shuffle-241517/webhook@sha256:654a5031789062f33b03d1a1004189895efe9df19fe762c688dff72522ce1a67"

	// Define the service
	// Define the service to deploy..
	// Wtf even is this
	service := &cloudrun.Service{
		ApiVersion: "v1alpha1",
		Kind:       "Service",
		Metadata: &cloudrun.ObjectMeta{
			Name:      "webhook3",
			Namespace: "default",
		},
		Spec: &cloudrun.ServiceSpec{
			RunLatest: &cloudrun.ServiceSpecRunLatest{
				Configuration: &cloudrun.ConfigurationSpec{
					RevisionTemplate: &cloudrun.RevisionTemplate{
						Metadata: &cloudrun.ObjectMeta{
							Name:         "Helo",
							GenerateName: "1",
						},
						Spec: &cloudrun.RevisionSpec{
							Container: &cloudrun.Container{
								Image: imagename,
								Name:  "webhook",
							},
							Containers: []*cloudrun.Container{
								&cloudrun.Container{
									Image: imagename,
									Name:  "webhook",
								},
							},
							ContainerConcurrency: 80,
							ServingState:         "ACTIVE",
							TimeoutSeconds:       300,
						},
					},
				},
			},
		},
	}

	// Deploy the previously described service to all locations
	// Locations are the same as "parent" in other API calls, AKA:
	// projects/{projectname}/locations/{locationName}
	for _, location := range allLocations {
		projectsLocationsServicesCreateCall := projectsLocationsService.Services.Create(location, service)
		service, err = projectsLocationsServicesCreateCall.Do()
		log.Println(service, err)
		if err != nil {
			log.Fatalf("Error creating new locationservice: %s", err)
		}
	}

	//log.Printf("%#v", projectsLocationsService)
	//projectsLocationsServicesCreateCall := projectsLocationsService.Services.Create(parent, service)
	//log.Printf("%#v", projectsLocationsServicesCreateCall)

	//service, err = projectsLocationsServicesCreateCall.Do()
	//if err != nil {
	//	log.Println(err)
	//}

	log.Println(service)
}
