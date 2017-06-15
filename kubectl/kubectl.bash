#!/bin/bash

# If there is no current context, get one.
if [[ $(kubectl config current-context 2> /dev/null) == "" ]]; then
	cluster=$(gcloud config get-value container/cluster 2> /dev/null)
	zone=$(gcloud config get-value compute/zone 2> /dev/null)
	function var_usage() {
		cat <<EOF
No cluster is set. To set the cluster (and the zone where it is found), set the environment variables
  CLOUDSDK_COMPUTE_ZONE=<cluster zone>
  CLOUDSDK_CONTAINER_CLUSTER=<cluster name>
EOF
		exit 1
	}

	[[ -z "$cluster" ]] && var_usage
	[[ -z "$zone" ]] && var_usage
	err=$(mktemp)
	gcloud container clusters get-credentials --zone "$zone" "$cluster" 2> "$err"
	if [[ "$?" != 0 ]]; then
		cat "$err"
		exit 1
	fi
fi

kubectl "$@"
