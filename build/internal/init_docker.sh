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
TERRAFORM_VERSION="0.13.1"
PROVIDER_BRANCH_NAME="develop"

STASH_SERVER=git.source.akamai.com
GIT_IP=$(dig +short $STASH_SERVER)
[ -z "$GIT_IP" ] && echo "Aborting - Can not reach $STASH_SERVER. Check VPN" && exit 1 || echo "Resolved $STASH_SERVER preparing build"

read -p "Pull latest changes and rebuild image? (Enter yes(default) or no):" buildflag
if [ "$buildflag" = "no" ];
then
   PROVIDER_BRANCH_HASH="$(git rev-parse --short HEAD)"
   PROVIDER_BRANCH_NAME="$(git rev-parse --abbrev-ref HEAD)"
   echo "Running existing local build from branch $PROVIDER_BRANCH_NAME at $PROVIDER_BRANCH_HASH"
else
     docker rm -f akatf-container 2> /dev/null || true
     docker image rm -f terraform/akamai:$PROVIDER_BRANCH_NAME 2> /dev/null || true
     PROVIDER_BRANCH_HASH="$(git rev-parse --short HEAD)"
     PROVIDER_BRANCH_NAME="$(git rev-parse --abbrev-ref HEAD)"
     read -p "Manually specify branch? if no then develop branch will be used. (Enter yes or no(default)):" branchflag
     if [ "$branchflag" = "yes" ];
        then
          read -p "Enter branch name without spaces:" branchname
          PROVIDER_BRANCH_NAME="${branchname}"
          echo "Running existing local build from branch $PROVIDER_BRANCH_NAME"
     else
        echo "defaulting to develop branch"
     fi

   DOCKER_BUILDKIT=1 docker build \
      -f build/internal/package/Dockerfile.qa \
      --build-arg TERRAFORM_VERSION=${TERRAFORM_VERSION} \
      --build-arg PROVIDER_BRANCH_NAME=${PROVIDER_BRANCH_NAME} \
      --ssh default \
      --no-cache \
      -t terraform/akamai:$PROVIDER_BRANCH_NAME .
fi

docker rm -f akatf-container 2> /dev/null || true
docker run -d -it --name akatf-container --entrypoint "/usr/bin/tail" \
        -e TF_LOG=DEBUG \
        -e TF_LOG_PATH="provider.log" \
        -v $HOME/.edgerc:/root/.edgerc:ro \
        -v $WORKDIR:/tf:rw \
        -w /tf/ \
        terraform/akamai:$PROVIDER_BRANCH_NAME -f /dev/null
docker exec -it akatf-container sh
docker rm -f akatf-container 2> /dev/null || true