#!/bin/bash

NAMESPACE="$1"

if [[ -z "${NAMESPACE}" ]]; then
  >&2 echo "Please pass namespace to delete"
  exit 1
else
  kubectl delete namespace "${NAMESPACE}"
fi
