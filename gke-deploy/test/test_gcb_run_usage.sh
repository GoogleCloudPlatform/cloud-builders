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

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)
cd "${SCRIPT_DIR}"

# Execute and Verify
gcloud builds submit --config cloudbuild_gcb_run_usage.yaml . --project="${GKE_DEPLOY_PROJECT}" \
|| fail "gcb build failed"

# Clean up

cd "${SCRIPT_DIR}"
./clean_cluster.sh "${NAMESPACE}"

echo -e
echo "Success!"
