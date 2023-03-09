#!/usr/bin/env bash

set +x
set -e

TAG="${DOCKER_REGISTRY}/${CI_PROJECT_NAME}:${CI_COMMIT_REF_SLUG}-${CI_PIPELINE_ID}"

docker build --no-cache --pull \
  --build-arg DOCKER_PROXY="${DOCKER_PROXY}" \
  --build-arg GITLAB_USER="${GITLAB_USER}" \
  --build-arg GITLAB_TOKEN="${GITLAB_TOKEN}" \
  -t "${TAG}" .
docker push "${TAG}"
