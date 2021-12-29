#!/bin/bash

set +x
set -e

TAG="${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${CI_COMMIT_REF_SLUG}-${CI_PIPELINE_ID}"

for host in ${DEPLOY_HOSTS}; do
        echo --- Deploying to "${host}" ---;
        export DOCKER_HOST=tcp://${host}:2375;
        echo Pulling image "${TAG}"...;
        pulled_image=$(docker pull "${TAG}" -q);
        if [[ "${pulled_image}" == "${TAG}" ]]; then
          echo Stopping and removing existing "${CONTAINER_NAME}" container...;
          docker stop "${CONTAINER_NAME}" > /dev/null || true;
          docker rm -f "${CONTAINER_NAME}" > /dev/null || true;
          echo Starting new "${CI_PROJECT_NAME}" container on port "${HTTP_PORT}"...;
          docker run -d \
            --name "${CONTAINER_NAME}" \
            --dns "${NAMESERVER_0}" \
            --dns "${NAMESERVER_1}" \
            --restart unless-stopped \
            -p "${HTTP_PORT}":"${HTTP_PORT}" \
            -e ELASTIC_SERVER \
            -e ELASTIC_USER \
            -e ELASTIC_PASSWORD \
            -e SQL_SERVER_MALIBU \
            -e SQL_PORT_MALIBU \
            -e SQL_USER_MALIBU \
            -e SQL_PASSWORD_MALIBU \
            -e SQL_DATABASE_MALIBU \
            -e MINIO_URL \
            -e MINIO_FOLDER_NAME \
            -e MINIO_BUCKET_NAME \
            -e MINIO_SECRET_ACCESS_KEY \
            -e MINIO_ACCESS_KEY_ID \
            -e MAX_COUNT_PARTS \
            -e SITE_BASE_HOST \
            -e BASE_SUB_FOLDER \
            -e CRON_PATTERN_MAKER \
            -e HTTP_PORT \
            -e RATE_LIMIT_EVERY \
            -e RATE_LIMIT_BURST \
            -e IS_DEVELOPMENT \
            -e FLUENTD_SIDECAR_ADDRESS \
            "${TAG}" > /dev/null;
          echo "Successfully deployed"
        else
          echo -e "Failed to pull image, skipping this host" >&2;
        fi;
done;
