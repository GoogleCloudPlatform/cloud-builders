#!/bin/bash

cd examples

if [[ ! -d "${BUILDER_OUTPUT}" ]] ; then
  export BUILDER_OUTPUT="$(mkdir builderoutput)"
fi

extract_output_id () {
  HEX="[[:xdigit:]]"
  UUID_PATTERN="$HEX{8}-$HEX{4}-$HEX{4}-$HEX{4}-$HEX{12}"
  grep -oE $UUID_PATTERN "${BUILDER_OUTPUT}/output"
}

bazel run --spawn_strategy=standalone //:checkargs --verbose_failures -- a b d
if [[ $? -eq 0 ]] ; then
  echo "Wrapper script returned successful for failed bazel run"
  exit 1
fi
INVOCATION_ID_1="$(extract_output_id)"
echo "Invocation ID #1: ${INVOCATION_ID_1}"

if [[ -z "${INVOCATION_ID_1}" ]] ; then
  echo "\$BUILDER_OUTPUT/output not written to after failed bazel run."
  exit 1
fi


bazel run --spawn_strategy=standalone //:checkargs --verbose_failures -- a b c
if [[ $? -ne 0 ]] ; then
  echo "Wrapper script returned failure for successful bazel run"
  exit 1
fi
INVOCATION_ID_2="$(extract_output_id)"
echo "Invocation ID #2: ${INVOCATION_ID_2}"

if [[ "${INVOCATION_ID_1}" = "${INVOCATION_ID_2}" ]] ; then
  echo "\$BUILDER_OUTPUT/output not written to after successful bazel run."
  exit 1
fi
