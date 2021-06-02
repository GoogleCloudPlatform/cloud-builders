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

NAMESPACE="test-local-server-dry-run"
OUTPUT="/var/tmp/gke-deploy-test/test_local_server_dry_run"

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)
cd "${SCRIPT_DIR}"

./clean_cluster.sh "${NAMESPACE}" || true  # Don't exit if this fails
rm -rf "${OUTPUT}"

# Create namespace yaml in temp dir using template
TEMP_NAMESPACE_DIR="${OUTPUT}/namespace"
mkdir -p "${TEMP_NAMESPACE_DIR}"
sed "s/@NAME@/${NAMESPACE}/g" namespace/namespace.yaml > "${TEMP_NAMESPACE_DIR}/namespace.yaml"

# Execute

gke-deploy run \
-f "${TEMP_NAMESPACE_DIR}" \
-p "${GKE_DEPLOY_PROJECT}" \
-c "${GKE_DEPLOY_CLUSTER}" \
-l "${GKE_DEPLOY_LOCATION}" \
-o "${OUTPUT}" \
-D \
|| fail "gke-deploy run failed"

# Verify

cd "${OUTPUT}"/expanded
grep -Fq "kind: Namespace" * || fail "Expanded Namespace was not created in ${OUTPUT}"

mkdir "${OUTPUT}"/check && cd "${OUTPUT}"/check
gcloud container clusters get-credentials "${GKE_DEPLOY_CLUSTER}" --zone "${GKE_DEPLOY_LOCATION}" --project "${GKE_DEPLOY_PROJECT}"
! kubectl get namespace "${NAMESPACE}" >/dev/null 2>&1 || fail "Dry run should not have created namespace."

echo -e
echo "Success!"
