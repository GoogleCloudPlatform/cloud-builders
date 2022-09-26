# Tool builder: `gcr.io/cloud-builders/twine`

The `gcr.io/cloud-builders/twine` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of Twine. We also do
not provide historical pinned versions of Twine.

While there is, as of this writing, no official Twine image, users may install
Twine from a Python image. The Python team provides `python` images that
support multiple tagged versions of Python on multiple platforms. Please visit
https://hub.docker.com/_/python for more details. The following example
illustrates usage of an official `python` image.

## Example:

```yaml
steps:
- name: 'python:slim'
  entrypoint: '/bin/sh'
  args:
    - -c
    - |
      python -m pip install --user twine keyrings.google-artifactregistry-auth
- name: 'python:slim'
  entrypoint: '/bin/sh'
  args:
    - -c
    - |
      python -m twine upload --verbose --repository-url 'https://my-region-python.pkg.dev/my-project/my-repo' dist/*
```

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
