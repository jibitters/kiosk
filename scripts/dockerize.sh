#!/usr/bin/env bash

# Extracting the version number from the executable file.
VERSION=$(find kiosk-linux-* | sed -E 's/kiosk-linux-v//')

IMAGE_NAME=jibitters/kiosk
case $TRAVIS_BRANCH in
	master)
	    IMAGE_TAG=$VERSION
		;;
	develop)
		IMAGE_TAG=$VERSION-INTEGRATION
		;;
	*)
		IMAGE_TAG=$VERSION-DEVELOPMENT
		;;
esac

echo "About to build the $IMAGE_NAME:$IMAGE_TAG image"
docker build -t "$IMAGE_NAME:$IMAGE_TAG" .

echo "Signing into registry!"
docker login -u "$REGISTRY_USER" -p "$REGISTRY_PASSWORD"

echo "Pushing the $IMAGE_NAME:$IMAGE_TAG ..."
docker push "$IMAGE_NAME:$IMAGE_TAG"

echo "Removing the $IMAGE_NAME:$IMAGE_TAG from local build ..."
docker rmi -f "$IMAGE_NAME:$IMAGE_TAG"
