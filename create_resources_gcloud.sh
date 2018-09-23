#!/bin/bash

printf "\nOnce you've logged into your project, enter the project id below."
GCLOUD_PROJECT="$1"
if [ -z "$1" ]; then
echo 
read -p 'Please enter project id: ' GCLOUD_PROJECT
echo
fi

if [ -z "$GCLOUD_PROJECT" ];
then
  cat<<KEYERR
  ***WARNING: Inform a valid gcloud project id     
KEYERR
exit 1
fi

LOCATION="us-central1"
gcloud config set project $GCLOUD_PROJECT

# Enable Services
SERVICES=(
  cloudiot
  firebase
)

for svc in ${SERVICES[@]};
do
  printf "Enabling Service: ${svc}\n"
  gcloud services enable ${svc}.googleapis.com
done

# Create pubsub topic & save topic in env variable 
gcloud pubsub topics create pubsub-voice-robot-arm
gcloud pubsub subscriptions create writer-pubsub-voice-robot-arm --topic=pubsub-voice-robot-arm  --topic-project=$GCLOUD_PROJECT
gcloud iot registries create robot-arm-registry --region=$LOCATION --event-notification-config=topic=pubsub-voice-robot-arm

DEVICE_ID="$2"
if [ -z "$2" ]; then
echo 
read -p 'Please enter device id: ' DEVICE_ID
echo
fi

# Create keys 
openssl req -x509 -nodes -newkey rsa:2048 \
-keyout ${DEVICE_ID}.key.pem \
-out ${DEVICE_ID}.crt.pem \
-days 365 \
-subj "/CN=unused"

# Register devices
gcloud iot devices create $DEVICE_ID  \
--region=$LOCATION \
--registry=robot-arm-registry \
--project=$GCLOUD_PROJECT \
 --public-key path=${DEVICE_ID}.crt.pem,type=rsa-x509-pem