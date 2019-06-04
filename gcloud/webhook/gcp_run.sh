#!/bin/bash
projectname=""
echo "Running and building webhook for project $projectname"
docker build . -t gcr.io/$projectname/webhook
docker push gcr.io/$projectname/webhook

# Deploying
# gcloud beta run deploy webhook --image gcr.io/$projectname/webhook
