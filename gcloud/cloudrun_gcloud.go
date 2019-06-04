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

	// Defines a the projectname, the servicename to use and an existing image to use
	projectId := "shuffle-241517"
	imagename := "gcr.io/shuffle-241517/webhook:sha256:875380c5b746f028f73fb892c20d8b40e4eb3c2bde6914af07f8941948dd91ed"
	servicename := "webhook2"

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

	// Define the service to deploy
	// Wtf even is this
	// Metadata initializers
	// SOO MANY LAYERS OF BULLSHIT (:
	tmpservice := &cloudrun.Service{
		ApiVersion: "serving.knative.dev/v1alpha1",
		Kind:       "Service",
		Metadata: &cloudrun.ObjectMeta{
			Name:      servicename,
			Namespace: projectId,
		},
		Spec: &cloudrun.ServiceSpec{
			RunLatest: &cloudrun.ServiceSpecRunLatest{
				Configuration: &cloudrun.ConfigurationSpec{
					RevisionTemplate: &cloudrun.RevisionTemplate{
						Metadata: &cloudrun.ObjectMeta{
							DeletionGracePeriodSeconds: 0,
						},
						Spec: &cloudrun.RevisionSpec{
							Container: &cloudrun.Container{
								Image: imagename,
								Resources: &cloudrun.ResourceRequirements{
									Limits: map[string]string{"memory": "256Mi"},
								},
								Stdin:     false,
								StdinOnce: false,
								Tty:       false,
							},

							ContainerConcurrency: 80,
							TimeoutSeconds:       300,
						},
					},
				},
			},
		},
	}
	//Env: []*cloudrun.EnvVar{
	//								&cloudrun.EnvVar{
	//									Name:  "PORT",
	//									Value: "8080",
	//								},
	//							},

	// Deploy the previously described service to all locations
	// Locations are the same as "parent" in other API calls, AKA:
	// projects/{projectname}/locations/{locationName}
	for _, location := range allLocations {
		getService(projectsLocationsService, location)
		createService(projectsLocationsService, location, tmpservice)
	}
}

// Gets a service
func getService(projectsLocationsService *cloudrun.ProjectsLocationsService, location string) {
	projectsLocationsServicesGetCall := projectsLocationsService.Services.Get(fmt.Sprintf("%s/services/webhook", location))

	service, err := projectsLocationsServicesGetCall.Do()
	if err != nil {
		log.Println("Error creating new locationservice: %s", err)
	}

	_ = service
}

func createService(projectsLocationsService *cloudrun.ProjectsLocationsService, location string, service *cloudrun.Service) {
	projectsLocationsServicesCreateCall := projectsLocationsService.Services.Create(location, service)
	service, err := projectsLocationsServicesCreateCall.Do()
	log.Println(service, err)
	if err != nil {
		log.Fatalf("Error creating new locationservice: %s", err)
	}

	log.Printf("%#v", service.Spec)
}
