#!/bin/bash

set -e  # Fail if any command below fails

function fail() {
  echo -e
  echo "Failed: $1"
  exit 1
}

# Prepare

[ "${GKE_DEPLOY_PROJECT}" ] || fail "Please set GKE_DEPLOY_PROJECT"
[ "${GKE_DEPLOY_CLUSTER}" ] || fail "Please set GKE_DEPLOY_CLUSTER"
[ "${GKE_DEPLOY_LOCATION}" ] || fail "Please set GKE_DEPLOY_LOCATION"

NAMESPACE="test-local-run-all"
OUTPUT="/var/tmp/gke-deploy-test/test_local_run_all"

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)
cd "${SCRIPT_DIR}"

./clean_cluster.sh "${NAMESPACE}" || true  # Don't exit if this fails
rm -rf "${OUTPUT}"

# Execute

gke-deploy run \
-f configs-all \
-i gcr.io/google-containers/nginx \
-a "test-name" \
-v "test-version" \
-L "foo=bar" \
-A "hi=bye" \
-p "${GKE_DEPLOY_PROJECT}" \
-c "${GKE_DEPLOY_CLUSTER}" \
-l "${GKE_DEPLOY_LOCATION}" \
-n "${NAMESPACE}" \
-o "${OUTPUT}" \
|| fail "gke-deploy run failed"

# Clean up

cd "${SCRIPT_DIR}"
./clean_cluster.sh "${NAMESPACE}"

echo -e
echo "Success!"
