#!/bin/sh

set -e

VERSION="$1"
PROJECT_ID="flawless-haven-446608-d0"
REGION="asia-northeast3"
CLOUD_RUN_SERVICE_NAME="podoroot"
IMAGE_TAG="${REGION}-docker.pkg.dev/${PROJECT_ID}/podossaem/podoroot:${VERSION}"
IMAGE_TAG_LATEST="${REGION}-docker.pkg.dev/${PROJECT_ID}/podossaem/podoadmin:latest"

TOKEN_LEN=${#GITHUB_TOKEN}
if [ "$TOKEN_LEN" -gt 0 ]; then
  echo "[INFO] GITHUB_TOKEN is set. (length: $TOKEN_LEN)"
else
  echo "[ERROR] GITHUB_TOKEN is not set or empty."
  exit 1
fi

docker buildx use podo-builder

docker buildx build \
  --platform linux/amd64 \
  --build-arg PHASE=prod \
  --build-arg GTK=${GITHUB_TOKEN} \
  -t "${IMAGE_TAG}" \
  -t "${IMAGE_TAG_LATEST}" \
  --push \
  --target runner \
  --cache-to=type=inline \
  --cache-from=type=registry,ref="${IMAGE_TAG_LATEST}" \
  .

docker push "${IMAGE_TAG_LATEST}"

gcloud run deploy ${CLOUD_RUN_SERVICE_NAME} \
  --image ${IMAGE_TAG} \
  --platform managed \
  --region ${REGION} \
  --allow-unauthenticated
