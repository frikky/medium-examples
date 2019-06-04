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
	projectId := "shuffle-241517"

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

	// Metadata initializers
	tmpservice := &cloudrun.Service{
		ApiVersion: "serving.knative.dev/v1alpha1",
		Kind:       "Service",
		Metadata: &cloudrun.ObjectMeta{
			Name:            "webhook3",
			Namespace:       projectId,
			ResourceVersion: "AAWKf7cmgXg",
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
	_ = tmpservice

	// Deploy the previously described service to all locations
	// Locations are the same as "parent" in other API calls, AKA:
	// projects/{projectname}/locations/{locationName}
	for _, location := range allLocations {
		//projectsLocationsServicesGetCall := projectsLocationsService.Services.Get(fmt.Sprintf("%s/services/webhook", location))

		//service, err := projectsLocationsServicesGetCall.Do()
		//if err != nil {
		//	log.Fatalf("Error creating new locationservice: %s", err)
		//}

		////log.Printf("%#v", service.Metadata)
		//log.Printf("%#v", service.Spec.RunLatest.Configuration.RevisionTemplate.Metadata)
		//log.Printf("%#v", service.Spec.RunLatest.Configuration.RevisionTemplate.Spec.Container.Resources)

		projectsLocationsServicesCreateCall := projectsLocationsService.Services.Create(location, tmpservice)
		service, err := projectsLocationsServicesCreateCall.Do()
		log.Println(service, err)
		if err != nil {
			log.Fatalf("Error creating new locationservice: %s", err)
		}

		log.Printf("%#v", service.Spec)
	}
}
