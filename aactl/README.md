# aactl

Google Artifact Analysis (AA) service data import utility, supports OSS vulnerability scanner reports, SLSA provenance, and sigstore attestations.

AACTL is a tool that allows Google Cloud customers who use Artifact Registry to ingest vulnerabilities detected by supported scanning tools. Once ingested, vulnerabilities will be stored & managed alongside vulnerabilities detected by Artifact Analysis. Vulnerabilities are viewable within Artifact Registry UI, SDS Security Insights, gcloud, and Artifact Analysis API (Container Analysis).

AACTL can also ingest SLSA Build Provenance generated by SLSA GitHub Generator.

## Adding aactl step to your pipeline

To use `aactl` in your GCB pipeline you will need to add the following step:

```yaml
- id: import
  name: gcr.io/$PROJECT_ID/aactl
  waitFor: [scan]
  args: ['import', '--project', '$PROJECT_ID', '--source', '${_IMAGE}', '--file', 'report.json', '--format', 'grype']
```

See [examples/cloudbuild.yaml](examples/cloudbuild.yaml) for full example scanner (using `grype`), and vulnerability import.


## Verifying imported vulnerabilities

To review the imported vulnerabilities in GCP:

```shell
gcloud artifacts docker images list $repo \
  --show-occurrences \
  --format json \
  --occurrence-filter "kind=\"VULNERABILITY\" AND resource_url=\"https://$image\""
```
