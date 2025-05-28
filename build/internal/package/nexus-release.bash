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

  if [[ -d "test" ]]; then
    rm -rf test
  fi;

}

find_branch() {
  CURRENT_BRANCH=$GIT_BRANCH
  if [[ "$CURRENT_BRANCH" == "develop" ]]; then
    EDGEGRID_BRANCH="develop"
  else
    # find parent branch from which this branch was created, iterate over the list of branches from the history of commits
    branches=($(git log --pretty=format:'%D' | sed 's@HEAD -> @@' | grep . | sed 's@origin/@@g' | sed 's@release/.*@@g' | sed -E $'s@master, (.+)@\\1, master@g' | tr ', ' '\n' | grep -v 'tag:' | sed -E 's@^v([0-9]+\.?){2,}(-rc\.[0-9]+)?@@g' | grep -v release/ | grep -v HEAD | sed '/^$/d'))
    branches+=("develop") # guard to fallback to safe value if less branches than 5
    for branch in ${branches[*]}
    do
      echo "Checking branch '${branch}'"
      EDGEGRID_BRANCH=$branch

      if [[ "$index" -eq "5" ]]; then
        echo "Exceeding limit of checks, fallback to default branch 'develop'"
        EDGEGRID_BRANCH="develop"
        break
      fi
      index=$((index + 1))

      git -C ./akamaiopen-edgegrid-golang branch -r | grep $EDGEGRID_BRANCH > /dev/null
      if [[ $? -eq 0 ]]; then
        echo "There is matching EdgeGrid branch '${EDGEGRID_BRANCH}'"
        break
      fi
    done
  fi
  echo "Current branch is '${CURRENT_BRANCH}', matching EdgeGrid branch is '${EDGEGRID_BRANCH}'"
}

clone_edgegrid() {
  if [ ! -d "./akamaiopen-edgegrid-golang" ]
  then
    echo "First time build, cloning the 'akamaiopen-edgegrid-golang' repo"
    git clone ssh://git@git.source.akamai.com:7999/devexp/akamaiopen-edgegrid-golang.git
  else
    echo "Repository 'akamaiopen-edgegrid-golang' already exists, so only cleaning and updating it"
    pushd akamaiopen-edgegrid-golang
    git reset --hard
    git fetch -p
    popd
  fi
}

checkout_edgegrid() {
  pushd akamaiopen-edgegrid-golang
  git checkout $EDGEGRID_BRANCH -f
  git reset --hard origin/$EDGEGRID_BRANCH
  git pull
  popd
}

adjust_edgegrid() {
  go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v11="./akamaiopen-edgegrid-golang"
  go mod tidy
  git diff
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

outputs=()
parse_arguments
clean
if [[ "$RELEASE_TYPE" == "snapshot" ]]; then
  clone_edgegrid
  find_branch
  checkout_edgegrid
  adjust_edgegrid
fi
if ! ./build/internal/docker_jenkins.bash "$CURRENT_BRANCH" "$EDGEGRID_BRANCH"; then
  exit 1
fi
build
nexus_push
