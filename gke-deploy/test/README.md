# gke-deploy Tests

Scripts used to execute functional tests.

## Initialize project for gke-deploy tests

### Set test project on local machine

Add to .bashrc
```bash
export GKE_DEPLOY_PROJECT="<project>"
export GKE_DEPLOY_CLUSTER="<cluster>"
export GKE_DEPLOY_LOCATION="<location>"
```

### Create test cluster

```bash
gcloud container clusters create $GKE_DEPLOY_CLUSTER --num-nodes 3 --zone $GKE_DEPLOY_LOCATION --project $GKE_DEPLOY_PROJECT
```

## Tests

### test_gcb_run_usage.sh

Tests `gke-deploy run --help` on GCB. This can be used to check that the test container has the "latest" changes.

### test_gcb_run.sh

Tests `gke-deploy run` on GCB.

### test_gcb_prepare_apply.sh

Tests `gke-deploy prepare` and `gke-deploy apply` on GCB.

### test_local_run.sh

Tests `gke-deploy run` locally.

### test_local_run_all.sh

Tests `gke-deploy run` locally with all Kubernetes config types that gke-deploy has custom readiness implementations for.

### test_local_prepare_apply.sh

Tests `gke-deploy prepare` and `gke-deploy apply` locally.

### test_local_configs_expose.sh

Tests `gke-deploy run` locally with no -f flag and with -x flag.

### test_local_run_diff_namespaces.sh

Tests `gke-deploy run` locally with no -n flag where each config has its own namespace.
