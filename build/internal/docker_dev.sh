#!/bin/sh
# This script will build the associated golang library and use it to build the provider
# this project directory must be named terraform-provider-akamai and the edgegrid client
# golang library must be checked out in a sibling directory named akamaiopen-edgegrid-golang,
#
# You can set WORKDIR env variable to point to your local examples directory if running locally
# this will then open up an interactive shell mounted to the working directory.  If not set the
# current directory is used instead,
#
# Lastly you can quit the shell by typing 'exit' and press return key. This will destroy the container
WORKDIR="${WORKDIR-$(pwd)}"
echo "WORKDIR is $WORKDIR"
TERRAFORM_VERSION="0.13.5"
PROVIDER_BRANCH_HASH="$(git rev-parse --short HEAD)"
PROVIDER_BRANCH_NAME="$(git rev-parse --abbrev-ref HEAD)"
eTAG="$(git describe --tags --always)"
echo "Making dev build on branch $PROVIDER_BRANCH_NAME at hash $PROVIDER_BRANCH_HASH with tag $eTAG"

STASH_SERVER=git.source.akamai.com
GIT_IP=$(dig +short $STASH_SERVER)
[ -z "$GIT_IP" ] && echo "Aborting - Can not reach $STASH_SERVER. Check VPN" && exit 1 || echo "Resolved $STASH_SERVER preparing build"

echo Building terraform-provider-akamai:$eTAG
# need to go up to resolve both akamai-terraform-provider and akamaiopen-edgegrid-golang project directories.
# but we don't want to include everything in the docker image.  This excludes everything but the two directories we want
ls -A -1 .. | grep -v terraform-provider-akamai | grep -v akamaiopen-edgegrid-golang > ../.dockerignore
cd ..
DOCKER_BUILDKIT=1 docker build --no-cache -f terraform-provider-akamai/build/internal/package/Dockerfile.dev \
      --build-arg TERRAFORM_VERSION=${TERRAFORM_VERSION} \
      --ssh default \
      --no-cache \
      -t terraform/akamai-local:$eTAG .
cd -

docker rm -f akatf-dev-container 2> /dev/null || true
docker run -d -it --name akatf-dev-container --entrypoint "/usr/bin/tail" \
        -e TF_LOG=DEBUG \
        -e TF_LOG_PATH="provider.log" \
        -v $HOME/.edgerc:/root/.edgerc:ro \
        -v $WORKDIR:/tf:rw \
        -w /tf/ \
        terraform/akamai-local:$eTAG -f /dev/null
docker exec -it akatf-dev-container sh
docker rm -f akatf-dev-container
echo Container destruction complete. Using Akamai provider plugin version:$eTAG
