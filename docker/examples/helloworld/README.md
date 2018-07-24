# Example

This directory contains a simple example that uses the docker build step.

To build this "Hello, world!" app, run:

```
gcloud builds submit --tag=gcr.io/my-project/sample-image .
```

Once your build is successful, you can run the app like this:

```
gcloud docker -- run --rm gcr.io/my-project/sample-image
Hello, world! The time is Wed Oct 19 17:02:44 UTC 2016.
```
