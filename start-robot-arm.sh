#!/bin/sh

export PROJECT_ID="voice-robot-arm"
export REGION="us-central1"
export REGISTRY_ID="robot-arm-registry"
export DEVICE_KEY="/home/pi/go/orangepizero.key.pem"
export CA_CERTS="/home/pi/go/roots.pem"
sudo /home/pi/go/robot-arm > /home/pi/go/out.txt 2>&1 