#!/bin/bash

set -ex

while read LINE; do \
  mvn --batch-mode org.apache.maven.plugins:maven-dependency-plugin:3.0.0:get -Dartifact=$LINE ; \
done < /builder/deps.txt
