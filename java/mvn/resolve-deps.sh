#!/bin/bash

set -ex

MAX_RETRIES=5

while read LINE
do
  # fetch the artifact, retrying up to 5 times
  i=0
  until [[ $i -ge $MAX_RETRIES ]]
  do
    mvn --batch-mode org.apache.maven.plugins:maven-dependency-plugin:3.0.0:get -Dartifact="$LINE" && break;
    i=$((i + 1))
    sleep 1
  done

  if [[ $i -ge $MAX_RETRIES ]]; then
    echo "Failed to fetch artifact \"$LINE\" after $MAX_RETRIES attempts"
    exit 1
  fi

done < /builder/deps.txt
