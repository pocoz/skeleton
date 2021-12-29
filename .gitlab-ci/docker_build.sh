#!/usr/bin/env bash

set +x
set -e

TAG="${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${CI_COMMIT_REF_SLUG}-${CI_PIPELINE_ID}"

docker build --no-cache --pull -t "${TAG}" .
docker push "${TAG}"
