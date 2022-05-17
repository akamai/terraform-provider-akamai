#!/usr/bin/env bash

# This script executes a "dry run" simulating what the goreleaser GitHub action does for every release
# It is executed on last commit from develop branch.

# Script will end immediately when some command exits with a non-zero exit code.
set -e

SSH_PRV_KEY="$(cat ~/.ssh/id_rsa)"
SSH_PUB_KEY="$(cat ~/.ssh/id_rsa.pub)"
SSH_KNOWN_HOSTS="$(cat ~/.ssh/known_hosts)"
IMAGE="tfp-releaser"
CONTAINER="${IMAGE}-container"

clean_up() {
  numConts=$( docker container ls | grep -w "${CONTAINER}" | wc -l )
  if [[ $numConts -ne 0 ]]; then
    docker container rm -f "${CONTAINER}"
    echo "removed $numConts containers"
  else
    echo "no containers to remove"
  fi

  numImages=$( docker image ls | grep -w "${IMAGE}" | wc -l )
  if [[ $numImages -ne 0 ]]; then
    docker image rm -f "${IMAGE}"
    echo "removed $numImages images"
  else
    echo "no images to remove"
  fi
  return 0
}

clean_up

echo "building image:"
cd build/internal/releaser
docker build -f Dockerfile -t "${IMAGE}" .

echo "starting container:"
docker container run --name "${CONTAINER}" -d -it \
  -e SSH_PUB_KEY="${SSH_PUB_KEY}" \
  -e SSH_PRV_KEY="${SSH_PRV_KEY}" \
  -e SSH_KNOWN_HOSTS="${SSH_KNOWN_HOSTS}"  "${IMAGE}"

echo "cloning repositories:"
docker container exec "${CONTAINER}" bash -c clone_repos.bash

echo "executing goreleaser build:"
docker container exec "${CONTAINER}" bash -c goreleaser_build.bash

echo "smoke test:"
docker container exec ${CONTAINER} bash -c smoke_tests.bash

clean_up
