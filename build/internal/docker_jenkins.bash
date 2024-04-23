#!/usr/bin/env bash
# This script will build the provider and associated library after checking out from git on jenkins.
#
# It uses the same docker image for all builds unless RELOAD_DOCKER_IMAGE parameter is set true.

# Script will end immediately when some command exits with a non-zero exit code.
set -e

PROVIDER_BRANCH_NAME="${1:-develop}"
EDGEGRID_BRANCH_NAME="${2:-develop}"
RELOAD_DOCKER_IMAGE="${3:-false}"

# Recalculate DOCKER_IMAGE_SIZE if any changes to dockerfile.
TIMEOUT="40m"
DOCKER_IMAGE_SIZE="554443852"

SSH_PRV_KEY="$(cat ~/.ssh/id_rsa)"
SSH_PUB_KEY="$(cat ~/.ssh/id_rsa.pub)"
SSH_KNOWN_HOSTS="$(cat ~/.ssh/known_hosts)"
SSH_CONFIG="PubkeyAcceptedKeyTypes +ssh-rsa"

COVERAGE_DIR=test/coverage
COVERAGE_PROFILE="$COVERAGE_DIR"/profile.out
COVERAGE_XML="$COVERAGE_DIR"/coverage.xml
COVERAGE_HTML="$COVERAGE_DIR"/index.html

WORKDIR="${WORKDIR-$(pwd)}"
echo "WORKDIR is $WORKDIR"
TERRAFORM_VERSION="1.4.6"

STASH_SERVER=git.source.akamai.com
GIT_IP=$(dig +short $STASH_SERVER)
[ -z "$GIT_IP" ] && echo "Aborting - Can not reach $STASH_SERVER." && exit 1 || echo "Resolved $STASH_SERVER, preparing build"

eTAG="$(git describe --tags --always)"
PROVIDER_BRANCH_HASH="$(git rev-parse --short HEAD)"
echo "Making build on branch $PROVIDER_BRANCH_NAME at hash $PROVIDER_BRANCH_HASH with tag $eTAG"

mkdir -p $COVERAGE_DIR

docker rm -f akatf-container 2> /dev/null || true

# Remove docker image if RELOAD_DOCKER_IMAGE is true
if [[ "$RELOAD_DOCKER_IMAGE" == true ]]; then
  echo "Removing docker image terraform/akamai:terraform-provider-akamai if exists"
  docker image rm -f terraform/akamai:terraform-provider-akamai 2> /dev/null || true
fi

if [[ "$(docker images -q terraform/akamai:terraform-provider-akamai 2> /dev/null)" == "" ||
      "$(docker inspect -f '{{ .Size }}' terraform/akamai:terraform-provider-akamai)" != "$DOCKER_IMAGE_SIZE" ]]; then
  echo "Building new image terraform/akamai:terraform-provider-akamai"
  DOCKER_BUILDKIT=1 docker build \
    -f build/internal/package/Dockerfile \
    --build-arg TERRAFORM_VERSION=${TERRAFORM_VERSION} \
    --no-cache \
    -t terraform/akamai:terraform-provider-akamai .
fi

echo "Creating docker container"
docker run -d -it --name akatf-container --entrypoint "/usr/bin/tail" \
        -e TF_LOG=DEBUG \
        -e TF_LOG_PATH="provider.log" \
        -e COVERMODE="atomic" \
        -e EDGEGRID_BRANCH_NAME="$EDGEGRID_BRANCH_NAME" \
        -e PROVIDER_BRANCH_NAME="$PROVIDER_BRANCH_NAME" \
        -e SSH_PUB_KEY="${SSH_PUB_KEY}" \
        -e SSH_PRV_KEY="${SSH_PRV_KEY}" \
        -e SSH_KNOWN_HOSTS="${SSH_KNOWN_HOSTS}" \
        -e SSH_CONFIG="${SSH_CONFIG}" \
        -e TIMEOUT="$TIMEOUT" \
        -e TERRAFORM_VERSION="$TERRAFORM_VERSION" \
        -v "$HOME"/.ssh/id_rsa=/root/id_rsa \
        -v "$HOME"/.ssh/id_rsa.pub=/root/id_rsa.pub \
        -v "$HOME"/.ssh/known_hosts=/root/known_hosts \
        -v "$WORKDIR"/.edgerc:/root/.edgerc:ro \
        -w /tf/ \
        terraform/akamai:terraform-provider-akamai -f /dev/null

docker exec akatf-container sh -c 'echo "$SSH_KNOWN_HOSTS" > /root/.ssh/known_hosts;
                                   echo "$SSH_PUB_KEY" > /root/.ssh/id_rsa.pub;
                                   echo "$SSH_PRV_KEY" > /root/.ssh/id_rsa;
                                   echo "$SSH_CONFIG" > /root/.ssh/config;
                                   chmod 700 /root/.ssh;
                                   chmod 600 /root/.ssh/id_rsa;
                                   chmod 644 /root/.ssh/id_rsa.pub /root/.ssh/known_hosts /root/.ssh/config'

echo "Cloning repos"
docker exec akatf-container sh -c 'git clone ssh://git@git.source.akamai.com:7999/devexp/terraform-provider-akamai.git;
                                   git clone ssh://git@git.source.akamai.com:7999/devexp/akamaiopen-edgegrid-golang.git edgegrid'

echo "Checkout branches"
docker exec akatf-container sh -c 'cd edgegrid; git checkout ${EDGEGRID_BRANCH_NAME};
                                   cd ../terraform-provider-akamai; git checkout ${PROVIDER_BRANCH_NAME};
                                   go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v8=../edgegrid'

echo "Installing terraform"
docker exec akatf-container sh -c 'cd terraform-provider-akamai; make tools.terraform'

echo "Running go mod tidy"
docker exec akatf-container sh -c 'cd terraform-provider-akamai; make tidy'

echo "Running golangci-lint"
docker exec akatf-container sh -c 'cd terraform-provider-akamai; make lint'

echo "Running terraform fmt"
docker exec akatf-container sh -c 'cd terraform-provider-akamai; make terraform-fmtcheck'

echo "Running tflint on examples"
docker exec akatf-container sh -c 'cd terraform-provider-akamai; make terraform-lint'

echo "Running tests with xUnit output"
docker exec akatf-container sh -c 'cd terraform-provider-akamai;
                                   2>&1 go test -timeout $TIMEOUT -v -coverpkg=./... -coverprofile=../profile.out -covermode=$COVERMODE -skip TestClient_DefaultRetryPolicy_TLS ./... | tee ../tests.output'
docker exec akatf-container sh -c 'cat tests.output | go-junit-report' > test/tests.xml
docker exec akatf-container sh -c 'cat tests.output' > test/tests.output
sed -i -e 's/skip=/skipped=/g;s/ failures=/ errors="0" failures=/g' test/tests.xml

echo "Creating coverage files"
docker exec akatf-container sh -c 'cd terraform-provider-akamai; go tool cover -html=../profile.out -o ../index.html;
                                   gocov convert ../profile.out | gocov-xml > ../coverage.xml'
docker exec akatf-container sh -c 'cat profile.out' > "$COVERAGE_PROFILE"
docker exec akatf-container sh -c 'cat index.html' > "$COVERAGE_HTML"
docker exec akatf-container sh -c 'cat coverage.xml' > "$COVERAGE_XML"

echo "Creating docker build"
docker exec akatf-container sh -c 'cd terraform-provider-akamai; make build'

docker rm -f akatf-container 2> /dev/null || true
