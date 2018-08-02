## Google Cloud Build integration with JFrog Artifactory

There are times when you’ll want to use Cloud Builder with a private repository such as JFrog Artifactory, the universal artifact manger. Here are few use cases and there will certainly be more -
* Containerize an existing application that’s stored in a private repository?
* Containerize an application that relies on both private and public dependencies?
* Quicken the build and deployment time by caching dependencies?

This readme walks through the steps required to configure CloudBuild to work with JFrog Artifactory.

#### Step 1: Security

Since credentials are involved to authenticate with Artifactory, it is extremely important to ensure that credentials are passed in a secure manner in the cloudbuild.yaml file.

It is recommended to encrypt Artifactory API keys to make sure that encrypted credentials are used in Google Cloud Build. 

In order to do so, first create a Cloud KMS KeyRing and CryptoKey 


**How-to create KeyRing**

`gcloud kms keyrings create [KEYRING-NAME] --location=global
`


**How-to create CryptoKey**

`gcloud kms keys create [KEY-NAME] --location=global --keyring=[KEYRING-NAME] --purpose=encryption
`

Once the keyring and cryptokey are created, it can be used to encrypt strings and even a file that includes sensitive information.

**How-to encrypt API key**

`echo $RT_API_KEY | gcloud kms encrypt --plaintext-file=- --ciphertext-file=- --location=global --keyring=[KEYRING-NAME] --key=[KEY-NAME] | base64
`
This command will output an encrypted version of API KEY that will be referred as [ENCRYPTED_API_KEY] in the readme and sample scripts.


#### Step 2: Build a project with JFrog Artifactory as a source of truth for all types of binaries

`gcloud container builds submit --config=cloudbuild.yaml .`

NOTE: Make sure that the builder image exists before running the above step.

##### Four key steps that are part of the sample cloudbuild.yaml file:

* ###### Create a builder that includes JFrog CLI

JFrog CLI is package agnostic that means that the same version of CLI can be used to build maven, gradle, npm, Go, Conan, docker projects. 

This sample makes the base image configurable so that it's easy to generate a builder for a specific package type that also includes JFrog CLI.


```steps:
- name: 'gcr.io/cloud-builders/docker'
  args:
  - 'build'
  - '--build-arg=BASE_IMAGE=gcr.io/${PROJECT_ID}/mvn:3.3.9-jdk-8'
  - '--tag=gcr.io/$PROJECT_ID/java/jfrog:1.17.0'
  - '.'
  wait_for: ['-']
  
```

Once the builder is created, it can be used to build maven, gradle, npm, Go, docker based projects. 

NOTE: The example project builds a maven project, and hence relies on an existing builder image for maven projects. Make sure that the builder image exists before running the above step.
The builder image for maven is located at https://github.com/GoogleCloudPlatform/cloud-builders/blob/master/mvn/cloudbuild.yaml

* ###### Configure JFrog CLI to point to Jfrog Artifactory

```
- name: 'gcr.io/$PROJECT_ID/java/jfrog'
  entrypoint: 'bash'
  args: ['-c', 'jfrog rt c rt-mvn-repo --url=https://[ARTIFACTORY-URL]/artifactory --user=[ARTIFACTORY-USER] --password=$$APIKEY']
  secretEnv: ['APIKEY']
  dir: 'examples/maven-example'
```

**Note:** There is an added step in order to use the encrypted version of APIKEY
```
secrets:
- kmsKeyName: projects/[PROJECT]/locations/global/keyRings/[KEYRING-NAME]/cryptoKeys/[KEY-NAME]
  secretEnv:
    APIKEY: [ENCRYPTED_API_KEY]

```
* ###### Build a maven project
```
- name: 'gcr.io/$PROJECT_ID/java/jfrog'
  args: ['rt', 'mvn', "clean install", 'config.yaml', '--build-name=mybuild', '--build-number=$BUILD_ID']
  dir: 'examples/maven-example'
```
The step above refers to [config.yaml](./examples/maven-example/config.yaml) that specifies the maven repositories to use in JFrog Artifactory to pulland push snapshot and release maven artifacts. Additional information can be found [here](https://www.jfrog.com/confluence/display/CLI/CLI+for+JFrog+Artifactory#CLIforJFrogArtifactory-CreatingtheBuildConfigurationFile.1) 

**Note:** For other languages, it is recommended to use the corresponding BASE_IMAGE and follow this doc to build via JFrog CLI.


* ###### Containerize the app
```
- name: 'gcr.io/cloud-builders/docker'
  args:
  - 'build'
  - '--tag=gcr.io/$PROJECT_ID/java-app:${BUILD_ID}'
  - '.'
  dir: 'examples/maven-example'
  
```

Once the app is containerized, it can be deployed on GKE or any other compute target.
