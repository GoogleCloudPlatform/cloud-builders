#!/bin/bash

versions=( $(cd "$(gcloud info --format='value(installation.sdk_root)')/bin/" && ls kubectl.* | sed s/kubectl.//g) )

function var_usage() {
  cat <<EOF
  No cluster is set. To set the cluster (and the region/zone where it is found), set the environment variables
  CLOUDSDK_COMPUTE_REGION=<cluster region> (regional clusters)
  CLOUDSDK_COMPUTE_ZONE=<cluster zone>
  CLOUDSDK_CONTAINER_CLUSTER=<cluster name>

  Optionally, you can specify the kubectl version via KUBECTL_VERSION, taking one of the following values: ${versions[@]}
EOF

  exit 1
}

kubectl_cmd=kubectl${KUBECTL_VERSION:+.${KUBECTL_VERSION}}
cluster=${CLOUDSDK_CONTAINER_CLUSTER:-$(gcloud config get-value container/cluster 2> /dev/null)}
region=${CLOUDSDK_COMPUTE_REGION:-$(gcloud config get-value compute/region 2> /dev/null)}
zone=${CLOUDSDK_COMPUTE_ZONE:-$(gcloud config get-value compute/zone 2> /dev/null)}
project=${CLOUDSDK_CORE_PROJECT:-$(gcloud config get-value core/project 2> /dev/null)}

[[ -z "$cluster" ]] && var_usage
[ ! "$zone" -o "$region" ] && var_usage
if [[ -n "$KUBECTL_VERSION" ]] && [[ ! " ${versions[*]} "  =~ " ${KUBECTL_VERSION} "  ]]; then
  echo "Bad KUBECTL_VERSION \"${KUBECTL_VERSION}\". Expected one of ${versions[*]}" >&2
  exit 1
fi

if [ -n "$region" ]; then
  echo "Running: gcloud container clusters get-credentials --project=\"$project\" --region=\"$region\" \"$cluster\""
  gcloud container clusters get-credentials --project="$project" --region="$region" "$cluster" || exit
else
  echo "Running: gcloud container clusters get-credentials --project=\"$project\" --zone=\"$zone\" \"$cluster\""
  gcloud container clusters get-credentials --project="$project" --zone="$zone" "$cluster" || exit
 fi

echo "Running: ${kubectl_cmd}" "$@" >&2
exec "${kubectl_cmd}" "$@"
