# Hello World App

This example creates a simple web server that serves "Hello, world."

NOTE: The commands below assume that you have set `$PROJECT_ID` to the name of
your GCP project. If your project name is `my-project`, you can do this by
executing `export PROJECT_ID=my-project` prior to running the below commands.

First, build the app:

`gcloud builds submit --config=cloudbuild.yaml .`

If you have Docker installed locally, you can test your app. In one window, run:

`gcloud docker -- run -p 8080:8080 gcr.io/$PROJECT_ID/hello-app`

In a separate window, execute `curl localhost:8080` (or go there in your web
browser), and you should see your containerized web service respond `Hello,
world!`

If you don't have Docker installed locally, you can test your app by deploying
it into the Google App Engine flexible envionment:

`gcloud app deploy --image-url=gcr.io/$PROJECT_ID/hello-app app.yaml`

You can then hit the URL returned by that command to see the `Hello, world!` web
service response.

Don't forget to tear down your App Engine deployment to avoid billing charges
for a running service.
