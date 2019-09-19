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

OUTPUT="/var/tmp/gke-deploy-test/test_local_run_diff_namespaces"

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)
cd "${SCRIPT_DIR}"

./clean_cluster.sh a || true  # Don't exit if this fails
./clean_cluster.sh b || true  # Don't exit if this fails
./clean_cluster.sh c || true  # Don't exit if this fails
rm -rf "${OUTPUT}"

# Execute

gcloud container clusters get-credentials "${GKE_DEPLOY_CLUSTER}" --zone "${GKE_DEPLOY_LOCATION}" --project "${GKE_DEPLOY_PROJECT}"
kubectl create namespace a
kubectl create namespace b
kubectl create namespace c

gke-deploy run \
-f configs-diff-namespaces \
-i gcr.io/google-containers/nginx \
-a "test-name" \
-v "test-version" \
-L "foo=bar" \
-A "hi=bye" \
-p "${GKE_DEPLOY_PROJECT}" \
-c "${GKE_DEPLOY_CLUSTER}" \
-l "${GKE_DEPLOY_LOCATION}" \
-o "${OUTPUT}" \
|| fail "gke-deploy run failed"

# Verify

cd "${OUTPUT}"/expanded

mkdir "${OUTPUT}"/check && cd "${OUTPUT}"/check
gcloud container clusters get-credentials "${GKE_DEPLOY_CLUSTER}" --zone "${GKE_DEPLOY_LOCATION}" --project "${GKE_DEPLOY_PROJECT}"
kubectl get deployment test-deployment -n a -o yaml > deployment.yaml
grep -Fq "app.kubernetes.io/managed-by: gcp-cloud-build-deploy" deployment.yaml || fail "\"app.kubernetes.io/managed-by: gcp-cloud-build-deploy\" does not exist"
grep -Fq "app.kubernetes.io/name: test-name" deployment.yaml || fail "\"app.kubernetes.io/name: test-name\" label does not exist"
grep -Fq "app.kubernetes.io/version: test-version" deployment.yaml || fail "\"app.kubernetes.io/version: test-version\" label does note exist"
grep -Fq "foo: bar" deployment.yaml || fail "\"foo: bar\" label does not exist"
grep -Fq "hi: bye" deployment.yaml || fail "\"hi: bye\" annotation does not exist"
grep -Fq "gcr.io/google-containers/nginx@sha256" deployment.yaml || fail "\"gcr.io/google-containers/nginx@sha256\" container not found" # Can't guarantee digest won't change, but can check that a digest was added.
kubectl get hpa test-hpa -n b -o yaml
kubectl get service test-service -n c -o yaml

# Clean up

cd "${SCRIPT_DIR}"
./clean_cluster.sh a
./clean_cluster.sh b
./clean_cluster.sh c

echo -e
echo "Success!"
