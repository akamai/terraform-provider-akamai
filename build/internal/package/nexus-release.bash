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

find_branch() {
  # CURRENT_BRANCH=`git branch --show-current` # Note: --show-current requires Git 2.22.0+ which is not available on Jenkins agents
  CURRENT_BRANCH=$GIT_BRANCH
  if [[ "$CURRENT_BRANCH" == "develop" ]]; then
    EDGEGRID_BRANCH="v2"
  elif [[ $CURRENT_BRANCH =~ .*/sp-.* ]]; then
    EDGEGRID_BRANCH=$CURRENT_BRANCH
  else
    # find parent branch from which this branch was created
    EDGEGRID_BRANCH=`git log --pretty=format:'%D' HEAD^ | grep 'origin/' | head -n1 | sed 's@origin/@@' | sed 's@,.*@@'`
  fi
  git -C ./akamaiopen-edgegrid-golang branch -r | grep $EDGEGRID_BRANCH
  if [[ $? -ne 0 ]]; then
    echo "There is no matching EdgeGrid branch $EDGEGRID_BRANCH, fallback to 'v2'"
    EDGEGRID_BRANCH='v2'
  fi
  if [[ "$EDGEGRID_BRANCH" == "develop" ]]; then
    EDGEGRID_BRANCH="v2"
  fi
  echo "Current branch is $CURRENT_BRANCH, matching EdgeGrid branch is $EDGEGRID_BRANCH"
}

clone_edgegrid() {
  git clone ssh://git@git.source.akamai.com:7999/devexp/akamaiopen-edgegrid-golang.git
}

checkout_edgegrid() {
  cd akamaiopen-edgegrid-golang
  git checkout $EDGEGRID_BRANCH
  git pull
  cd ..
  pwd
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
    $MVN_PATH -B deploy:deploy-file -DgroupId="$NEXUS_GROUP" -DartifactId="$PACKAGE_NAME" -Dversion="$versionName" -Dpackaging="${binary#"$PACKAGE_NAME"-}-bin" -Dfile=bin/"$binary" -DrepositoryId=nexus -Durl="$repo_url" -DgeneratePom=false
  done
}

mod_edit() {
  edgegrid_version=$(go list -m -json -versions github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 | python3 -c "import sys, json; print(json.load(sys.stdin)['Version'])")
  go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v2="stash.akamai.com/fee/akamaiopen-edgegrid-golang.git/v2@${edgegrid_version}"
}

outputs=()
parse_arguments
clean
if [[ "$RELEASE_TYPE" == "snapshot" ]]; then
  clone_edgegrid
  find_branch
  checkout_edgegrid
  go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v2="./akamaiopen-edgegrid-golang"
else
  mod_edit
fi
build
nexus_push
