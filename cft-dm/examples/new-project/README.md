# New Project Example

This example demonstrates a build that:

1.  Clones the Deployment Manager and Cloud Foundation Toolkit sample templates
2.  Invokes the `cft` tool to deploy a project using Depoyment Manager

To run this example, make sure you have created a [`DM Creation Project`](https://github.com/GoogleCloudPlatform/deploymentmanager-samples/blob/master/community/cloud-foundation/templates/project/README.md). Ensure you have authenticated with gcloud, have configured permissions appropriately for cft
and run:

```
gcloud builds submit --config cloudbuild.yaml --substitutions=_CLOUD_FOUNDATION_PROJECT_ID="<DM Creation Project ID>",_CFT_ORGANIZATION_FOLDER_ID="<Folder ID of the Parent Folder for the New Project>",_CFT_BILLING_ACCOUNT_ID="<Billing Account ID to Attach the New Project to>",_CFT_CHILD_PROJECT="<Name of the New Project>"
```

The new Project will be created using Deployment Manager, with the Deployment located in the `DM Creation Project`
