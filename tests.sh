#! /bin/bash

GENERATE_LOGIN='{
  "detectIntentResponseId": "9afd8672-ae55-4df9-934a-50e5cc379928",
  "intentInfo": {
    "lastMatchedIntent": "projects/sodium-pathway-343117/locations/us-central1/agents/e2c8f950-cd0e-4a49-83ad-7d24d216f69e/intents/00000000-0000-0000-0000-000000000000",
    "displayName": "Default Welcome Intent",
    "confidence": 1.0
  },
  "pageInfo": {
    "currentPage": "projects/sodium-pathway-343117/locations/us-central1/agents/e2c8f950-cd0e-4a49-83ad-7d24d216f69e/flows/00000000-0000-0000-0000-000000000000/pages/START_PAGE",
    "displayName": "Start Page"
  },
  "sessionInfo": {
    "session": "projects/sodium-pathway-343117/locations/us-central1/agents/e2c8f950-cd0e-4a49-83ad-7d24d216f69e/sessions/27aa15-a61-c4b-dc3-34c78551d"
  },
  "fulfillmentInfo": {
    "tag": "generate-login"
  },
  "text": "hi",
  "languageCode": "en"
}'

if [ "$1" = "" ]
then 
    URL=http://localhost:8081/generate-login
else
    URL="$1"
fi

curl \
    -H 'Content-Type: application/json' \
    -d "$GENERATE_LOGIN" \
    $URL