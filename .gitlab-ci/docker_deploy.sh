#!/bin/bash

set +x
set -e

TAG="${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${CI_COMMIT_REF_SLUG}-${CI_PIPELINE_ID}"

for host in ${DEPLOY_HOST_VR}; do
        echo --- Deploying to ${host} ---;
        export DOCKER_HOST=tcp://${host}:2375;
        echo Pulling image "${TAG}"...;
        pulled_image=$(docker pull "${TAG}" -q);
        if [[ "${pulled_image}" == "${TAG}" ]]; then
          echo Stopping and removing existing "${CONTAINER_NAME}" container...;
          docker stop "${CONTAINER_NAME}" > /dev/null || true;
          docker rm -f "${CONTAINER_NAME}" > /dev/null || true;
          echo Starting new "${CI_PROJECT_NAME}" container without port...;
          docker run -d \
            --name "${CONTAINER_NAME}" \
            --log-driver fluentd --log-opt fluentd-address="${FLUENTD_ADDRESS}":"${FLUENTD_HAPROXY_PORT}" --log-opt tag="${FLUENTD_TAG}" --log-opt mode=non-blocking --log-opt max-buffer-size=300m --log-opt fluentd-sub-second-precision="true" \
            --dns "${NAMESERVER_0}" \
            --dns "${NAMESERVER_1}" \
            --restart unless-stopped \
            -e ELASTIC_SERVER \
            -e ELASTIC_USER \
            -e ELASTIC_PASSWORD \
            -e SQL_SERVER \
            -e SQL_PORT \
            -e SQL_USER \
            -e SQL_PASSWORD \
            -e SQL_DATABASE \
            --restart unless-stopped \
            "${TAG}" > /dev/null;
          echo "Successfully deployed"
        else
          echo -e "Failed to pull image, skipping this host" >&2;
        fi;
done;
