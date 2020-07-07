#!/bin/bash

set -euo pipefail
# Opinionated script to release on GitHub.
# This script runs in CircleCI, in a golang docker container from a folder that is a git repo.
# The script expects the binaries to reside in the build folder.

export TAG="`git tag --list --sort=v:refname --points-at HEAD | tail -1`"

echo "TAG=$TAG"
if [[ -z "$TAG" ]]; then
  echo "No tag to build."
  exit 1
fi

# Create GitHub release draft
draftdata="
{
  \"tag_name\": \"$TAG\",
  \"target_commitish\": \"master\",
  \"name\": \"$TAG\",
  \"body\": \"\",
  \"draft\": true,
  \"prerelease\": false
}
"
curl -s -S -X POST -u "${GITHUB_USERNAME}:${GITHUB_TOKEN}" https://api.github.com/repos/freshautomations/sconfig/releases --user-agent freshautomations -H "Accept: application/vnd.github.v3.json" -d "$draftdata" > draft.json
ERR=$?
if [[ $ERR -ne 0 ]]; then
  echo "ERROR: curl error, exitcode $ERR."
  exit $ERR
fi

export id="`jq -r .id draft.json`"
if [ -z "$id" ]; then
  echo "ERROR: Could not get draft id."
  exit 1
fi

echo "Release ID: ${id}"

# Upload binaries

for binary in sconfig_darwin_386 sconfig_darwin_amd64 sconfig_linux_386 sconfig_linux_amd64 sconfig_windows_386 sconfig_windows_amd64
do
echo -ne "Processing ${binary}... "
if [[ ! -f "build/${binary}" ]]; then
  echo "${binary} does not exist."
  continue
fi
curl -s -S -X POST -u "${GITHUB_USERNAME}:${GITHUB_TOKEN}" "https://uploads.github.com/repos/freshautomations/sconfig/releases/${id}/assets?name=${binary}" --user-agent freshautomations -H "Accept: application/vnd.github.v3.raw+json" -H "Content-Type: application/octet-stream" -H "Content-Encoding: utf8" --data-binary "@build/${binary}" > upload.json
ERR=$?
if [[ $ERR -ne 0 ]]; then
  echo "ERROR: curl error, exitcode $ERR."
  exit $ERR
fi

export uid="`jq -r .id upload.json`"
if [ -z "$uid" ]; then
  echo "ERROR: Could not get upload id for binary ${binary}."
  exit 1
fi

echo "uploaded binary ${binary}, id ${uid}."
done

rm draft.json
rm upload.json

# Publish release
releasedata="
{
  \"draft\": false,
  \"tag_name\": \"$TAG\"
}
"
curl -s -S -X POST -u "${GITHUB_USERNAME}:${GITHUB_TOKEN}" "https://api.github.com/repos/freshautomations/sconfig/releases/${id}" --user-agent script -H "Accept: application/vnd.github.v3.json" -d "$releasedata"
ERR=$?
if [[ $ERR -ne 0 ]]; then
  echo "ERROR: curl error, exitcode $ERR."
  exit $ERR
fi

