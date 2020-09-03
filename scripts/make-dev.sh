#!/usr/bin/env bash
# Script used to build and tag develop branch builds.

BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$BRANCH" != "develop" ]]; then
  echo 'Aborting script - only run on branch develop';
  exit 1;
fi

# run semtag even if called from another script
WORKDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TAG=$( ${WORKDIR}/semtag getcurrent);	   # get tag to describe current release

git tag $TAG;				           # uses tag set by previous command
go mod tidy;                           # better to use command below to run checks once code is clean
#/usr/bin/make --directory ${WORKDIR}/.. check build; # builds local project running all checks before building
/usr/bin/make --directory ${WORKDIR}/..; # builds local project
git push $TAG origin;                  # push tag to origin

