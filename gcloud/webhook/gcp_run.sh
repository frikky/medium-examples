#!/bin/bash
projectname="shuffle-241517"
echo "Running and building webhook for project $projectname"
docker build . -t gcr.io/$projectname/webhook
docker push gcr.io/$projectname/webhook

# Deploying
# gcloud beta run deploy webhook --image gcr.io/shuffle-241517/webhook
