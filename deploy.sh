#!/bin/sh

set -e

VERSION="$1"
PROJECT_ID="sbec-2025"
REGION="asia-northeast3"
CLOUD_RUN_SERVICE_NAME="sbec"
IMAGE_TAG="${REGION}-docker.pkg.dev/${PROJECT_ID}/sbec-docker-images/sbec:${VERSION}"

docker buildx use desktop-linux

docker buildx build \
  --platform linux/amd64 \
  -t "${IMAGE_TAG}" \
  --push \
  --target runner \
  .

# gcloud run deploy ${CLOUD_RUN_SERVICE_NAME} \
#   --image ${IMAGE_TAG} \
#   --platform managed \
#   --region ${REGION} \
#   --allow-unauthenticated
