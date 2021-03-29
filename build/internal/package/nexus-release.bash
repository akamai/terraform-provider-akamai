#!/usr/bin/env bash

PACKAGE_NAME=$1
RELEASE_TYPE=$2
MVN_PATH="/var/lib/jenkins/tools/hudson.tasks.Maven_MavenInstallation/Maven_3.3.9/bin/mvn"


######## EDIT HERE #########
NEXUS_GROUP=com.akamai.portalrel.test
PLATFORMS=("darwin/amd64" "linux/amd64")
# take version however you wish, ensure non-release builds have SNAPSHOT in name
VERSION="1.0.0"
############################


parse_arguments() {
  if [[ -z "$RELEASE_TYPE" || ! "$RELEASE_TYPE" =~ ^(snapshot|release)$ || -z "$PACKAGE_NAME" ]]; then
    echo "usage: $0 <package-name> <snapshot|release>"
    exit 1
  fi
}

clean() {
  if [[ -d "bin" ]]; then
    rm -rf bin
  fi;
}

build() {
  for platform in "${PLATFORMS[@]}"
  do
    IFS='/' read -r -a parsed <<< "$platform"
    GOOS=${parsed[0]}
    GOARCH=${parsed[1]}
    output_name="$PACKAGE_NAME-$GOOS-$GOARCH"
    echo "building $output_name..."
    if ! env GOOS="$GOOS" GOARCH="$GOARCH" go build -o bin/"$output_name"; then
        echo 'Error building version'
        exit 1
    fi
    outputs+=("$output_name")
  done
  echo "${outputs[0]}"
}

nexus_push() {
  sha=$(git rev-parse --short HEAD)
  if [[ $RELEASE_TYPE == "snapshot" ]]; then
    repo_url="https://lunabuild.akamai.com/nexus/content/repositories/snapshots"
    versionName="${VERSION}-${sha}-SNAPSHOT"
  else
    repo_url="https://lunabuild.akamai.com/nexus/content/repositories/releases"
    versionName="${VERSION}-${sha}"
  fi

  for binary in "${outputs[@]}"
  do
    $MVN_PATH deploy:deploy-file -DgroupId="$NEXUS_GROUP" -DartifactId="$PACKAGE_NAME" -Dversion="$versionName" -Dpackaging="${binary#"$PACKAGE_NAME"-}-bin" -Dfile=bin/"$binary" -DrepositoryId=nexus -Durl="$repo_url" -DgeneratePom=false
  done
}

mod_edit() {
  edgegrid_version=$(go list -m -json -versions github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 | python3 -c "import sys, json; print(json.load(sys.stdin)['Version'])")
  go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v2="stash.akamai.com/fee/akamaiopen-edgegrid-golang.git/v2@${edgegrid_version}"
}

outputs=()
clean
mod_edit
parse_arguments
build
nexus_push
