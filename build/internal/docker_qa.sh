#!/bin/sh
# This script will build the provider and associated library after checking out from git
#
# You can set WORKDIR env variable to point to your local examples directory if running locally
# this will then open up an interactive shell mounted to the working directory.  If not set the
# current directory is used instead,
#
# Lastly you can quit the shell by typing 'exit' and press return key. This will destroy the container
WORKDIR="${WORKDIR-$(pwd)}"
echo "WORKDIR is $WORKDIR"
TERRAFORM_VERSION="0.13.5"
PROVIDER_BRANCH_NAME="develop"

docker rm -f akatf-container 2> /dev/null || true
docker run -d -it --name akatf-container --entrypoint "/usr/bin/tail" \
        -e TF_LOG=DEBUG \
        -e TF_LOG_PATH="provider.log" \
        -v $HOME/.edgerc:/root/.edgerc:ro \
        -v $WORKDIR:/tf:rw \
        -w /tf/ \
        pdr.akamai.com/pulsar/terraform-provider:latest -f /dev/null
docker exec -it akatf-container sh
docker rm -f akatf-container 2> /dev/null || true
