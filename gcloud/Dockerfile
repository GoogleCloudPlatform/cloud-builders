FROM gcloud-slim

RUN apk update && \
    # JRE is required for cloud-datastore-emulator
    apk add --no-cache openjdk8-jre && \
    # Install all available components
    /google-cloud-sdk/bin/gcloud -q components install \
	alpha beta \
	app-engine-go \
	app-engine-java \
	app-engine-php \
	app-engine-python \
	app-engine-python-extras \
	bigtable \
	cbt \
	cloud-datastore-emulator \
	cloud-build-local \
	datalab \
	docker-credential-gcr \
	emulator-reverse-proxy \
	kubectl \
	pubsub-emulator \
	&& \
    # Update components
    /google-cloud-sdk/bin/gcloud -q components update && \
    /google-cloud-sdk/bin/gcloud components list
