#!/bin/bash

# always get a new context
cluster=$(gcloud config get-value container/cluster 2> /dev/null)
region=${CLOUDSDK_COMPUTE_REGION:-$(gcloud config get-value compute/region 2> /dev/null)}
zone=$(gcloud config get-value compute/zone 2> /dev/null)
project=$(gcloud config get-value core/project 2> /dev/null)

function var_usage() {
    cat <<EOF
No cluster is set. To set the cluster (and the region/zone where it is found), set the environment variables
CLOUDSDK_COMPUTE_REGION=<cluster region> (regional clusters)
CLOUDSDK_COMPUTE_ZONE=<cluster zone>
CLOUDSDK_CONTAINER_CLUSTER=<cluster name>
EOF
    exit 1
}

[[ -z "$cluster" ]] && var_usage
[ ! "$zone" -o "$region" ] && var_usage

if [ -n "$region" ]; then
  echo "Running: gcloud container clusters get-credentials --project=\"$project\" --region=\"$region\" \"$cluster\""
  gcloud container clusters get-credentials --project="$project" --region="$region" "$cluster" || exit
else
  echo "Running: gcloud container clusters get-credentials --project=\"$project\" --zone=\"$zone\" \"$cluster\""
  gcloud container clusters get-credentials --project="$project" --zone="$zone" "$cluster" || exit
fi

echo "Running: kubectl $@"
kubectl "$@"
