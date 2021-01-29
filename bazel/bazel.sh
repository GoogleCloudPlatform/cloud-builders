#!/bin/bash
set -e

# Always write the build_event_text_file to a temp file we create.
readonly BUILD_EVENT_FILE="$(mktemp --tmpdir build_event_text_file.XXXXX)"
readonly BUILD_EVENT_FILE_FLAG="--build_event_text_file=${BUILD_EVENT_FILE}"

# --build_event_text_file will be added to the command line as the last
# bazel-handled flag. This is not always the last argument due to
# the `bazel run` syntax for passing through arguments:
#
#  $ bazel [startup flags] run //:mytool [bazel run flags] -- [mytool args]
#
# To handle this syntax, we split the command line before the *first* "--"
# argument into two lists and put our new flag in the middle. When no "--"
# appears, this is equivalent to appending the new flag.
for ((i=1 ; i<=$# ; i++))
do
  [[ "${@:$i:1}" == "--" ]] && break
done
n=$(($#-$i+1))
# Insert our flag at index $i out of $n. Quoting is all required for expansion.
set -- "${@:1:$((i-1))}" "$BUILD_EVENT_FILE_FLAG" "${@:$i:$n}"

# Run bazel with the new command line and capture the exit code.
#
# This script has -e enabled, so if bazel exits non-zero, we need to both
# capture the exit code *and* replace it with zero. That is the purpose of the
# short-circuiting OR operator.
EXIT_CODE=0
/builder/bazel "$@" || EXIT_CODE=$?

# Parse out the UUID from the BEP output file and write it to
# $BUILDER_OUTPUT/output, whose first 4KB will be served inline in the Build.

# 'uuid' is a field in the BuildStarted event that is guaranteed to appear first
# in the build_event_text_file. That event might be larger than 4KB in text
# format though. This grep solution is crude, but allows us to avoid actually
# parsing the text proto file using real protobuf definitions.
readonly INVOCATION_ID_FIELD="$(grep -Eo 'uuid: \".{36}\"$' "${BUILD_EVENT_FILE}" | head -n 1)"
readonly INVOCATION_ID="${INVOCATION_ID_FIELD:7:-1}"
echo "{\"bazel.build/invocation_id\": \"${INVOCATION_ID}\"}" > "${BUILDER_OUTPUT}/output"

exit "${EXIT_CODE}"
