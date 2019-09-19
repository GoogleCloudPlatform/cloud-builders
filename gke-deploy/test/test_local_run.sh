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

NAMESPACE="test-local-run"
OUTPUT="/var/tmp/gke-deploy-test/test_local_run"

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)
cd "${SCRIPT_DIR}"

./clean_cluster.sh "${NAMESPACE}" || true  # Don't exit if this fails
rm -rf "${OUTPUT}"

# Execute

gke-deploy run \
-f configs \
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

# Verify

cd "${OUTPUT}"/expanded
[ -e namespace.yaml ] || fail "${OUTPUT}/expanded/namespace.yaml does not exist"

mkdir "${OUTPUT}"/check && cd "${OUTPUT}"/check
gcloud container clusters get-credentials "${GKE_DEPLOY_CLUSTER}" --zone "${GKE_DEPLOY_LOCATION}" --project "${GKE_DEPLOY_PROJECT}"
kubectl get deployment test-deployment -n "${NAMESPACE}" -o yaml > deployment.yaml
grep -Fq "app.kubernetes.io/managed-by: gcp-cloud-build-deploy" deployment.yaml || fail "\"app.kubernetes.io/managed-by: gcp-cloud-build-deploy\" does not exist"
grep -Fq "app.kubernetes.io/name: test-name" deployment.yaml || fail "\"app.kubernetes.io/name: test-name\" label does not exist"
grep -Fq "app.kubernetes.io/version: test-version" deployment.yaml || fail "\"app.kubernetes.io/version: test-version\" label does note exist"
grep -Fq "foo: bar" deployment.yaml || fail "\"foo: bar\" label does not exist"
grep -Fq "hi: bye" deployment.yaml || fail "\"hi: bye\" annotation does not exist"
grep -Fq "gcr.io/google-containers/nginx@sha256" deployment.yaml || fail "\"gcr.io/google-containers/nginx@sha256\" container not found" # Can't guarantee digest won't change, but can check that a digest was added.
kubectl get service test-service -n "${NAMESPACE}" -o yaml > service.yaml
SERVICE_IP="$(grep -F "ip:" service.yaml | awk '{print $NF}')" || fail "No IP in service"
curl "${SERVICE_IP}" || fail "Failed to curl service IP"

# Clean up

cd "${SCRIPT_DIR}"
./clean_cluster.sh "${NAMESPACE}"

echo -e
echo "Success!"
