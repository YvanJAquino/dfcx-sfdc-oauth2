#!/bin/bash

docker build -t gcr.io/holy-diver-297719/sfdc-oauth2 . 
docker push gcr.io/holy-diver-297719/sfdc-oauth2
gcloud alpha run deploy sfdc-oauth2 \
    --project holy-diver-297719 \
    --set-env-vars PROJECT_ID=holy-diver-297719 \
    --set-env-vars ENVIRONMENT=PRODUCTION \
    --image gcr.io/holy-diver-297719/sfdc-oauth2 \
    --timeout 5m \
    --no-cpu-throttling \
    --region us-east4 \
    --platform managed \
    --min-instances 0 \
    --max-instances 5 \
    --allow-unauthenticated


