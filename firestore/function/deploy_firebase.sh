gcloud functions deploy HandleUserChange --runtime go111 \
--trigger-event providers/cloud.firestore/eventTypes/document.update \
--trigger-resource "projects/YOUR-PROJECT/databases/(default)/documents/users/{userid}"
