#!/bin/bash
#
# Copyright 2016 Google, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

OLD_TAG="$1"
NEW_TAG="$2"

if [[ -z "$OLD_TAG" ]]; then
  echo "Error: existing docker image tag name must be provided as the 1st positional parameter." 1>&2
  exit 1
fi
if [[ -z "$NEW_TAG" ]]; then
  echo "Error: new docker image tag name must be provided as the 2nd positional parameter." 1>&2
  exit 1
fi

# TODO(jasmuth): Use the future GCR endpoint that allows remote retagging.

set -x
docker pull "$OLD_TAG"
docker tag "$OLD_TAG" "$NEW_TAG"
# We push here, instead of waiting for the worker to push later, to mimic
# the future behavior of this build step once it uses the GCR endpoint.
docker push "$NEW_TAG"
