/*
Copyright 2017 Google Inc. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main // import "gcp_test"

import (
	"log"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	v1compute "google.golang.org/api/compute/v1"
)

func main() {
	ctx := context.Background()

	projectID, err := metadata.ProjectID()
	if err != nil {
		log.Fatalf("Could not get Project ID: %v", err)
	}

	client := oauth2.NewClient(ctx, google.ComputeTokenSource(""))

	gce, err := v1compute.New(client)
	if err != nil {
		log.Fatalf("Could not create GCE client: %v", err)
	}

	resp, err := gce.Instances.AggregatedList(projectID).Do()
	if err != nil {
		log.Fatalf("Could not list instances: %v\nPerhaps you need to grant GCE viewer access to your container builder service account?", err)
	}

	for _, item := range resp.Items {
		for _, instance := range item.Instances {
			log.Print(instance.SelfLink)
		}
	}
}
