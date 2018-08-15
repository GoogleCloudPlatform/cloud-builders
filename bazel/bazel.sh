#!/bin/bash

set -e

# Always write the build_event_text_file to a temp file we create.
readonly BUILD_EVENT_FILE="$(mktemp --tmpdir build_event_text_file.XXXXX)"

/builder/bazel $@ --build_event_text_file="${BUILD_EVENT_FILE}"

# Parse out the UUID from the BEP output file and write it to
# $BUILDER_OUTPUT/output, whose first 4KB will be served inline in the Build.

# 'uuid' is a field in the BuildStarted event that is guaranteed to appear first
# in the build_event_text_file. That event might be larger than 4KB in text
# format though. This grep solution is crude, but allows us to avoid actually
# parsing the text proto file using real protobuf definitions.
readonly INVOCATION_ID_FIELD="$(grep -Eo 'uuid: \".{36}\"$' "${BUILD_EVENT_FILE}" | head -n 1)"
readonly INVOCATION_ID="${INVOCATION_ID_FIELD:7:-1}"
echo "{\"bazel.build/invocation_id\": \"${INVOCATION_ID}\"}" > "${BUILDER_OUTPUT}/output"
