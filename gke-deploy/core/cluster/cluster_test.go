/*
Copyright 2019 Google, Inc. All rights reserved.
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
package cluster

import (
	"context"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/resource"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/testservices"
)

func TestAuthorizeAccess(t *testing.T) {
	ctx := context.Background()
	clusterName := "test-cluster"
	clusterLocation := "us-east1-b"
	clusterProject := "my-project"
	gs := &testservices.TestGcloud{
		ContainerClustersGetCredentialsErr: nil,
	}

	if err := AuthorizeAccess(ctx, clusterName, clusterLocation, clusterProject, gs); err != nil {
		t.Errorf("AuthorizeAccess(ctx, %s, %s, gs) = %v; want <nil>", clusterName, clusterLocation, err)
	}
}

func TestAuthorizeAccessErrors(t *testing.T) {
	ctx := context.Background()
	clusterName := "test-cluster"
	clusterLocation := "us-east1-b"
	clusterProject := "my-project"
	gs := &testservices.TestGcloud{
		ContainerClustersGetCredentialsErr: fmt.Errorf("failed to get credentials of cluster"),
	}

	if err := AuthorizeAccess(ctx, clusterName, clusterLocation, clusterProject, gs); err == nil {
		t.Errorf("AuthorizeAccess(ctx, %s, %s, gs) = <nil>; want error", clusterName, clusterLocation)
	}
}

func TestApplyConfigFromString(t *testing.T) {
	configString := string(fileContents(t, "testing/deployment.yaml"))
	namespace := "default"
	ks := &testservices.TestKubectl{
		ApplyFromStringResponse: map[string][]error{
			configString: {nil},
		},
	}

	if err := ApplyConfigFromString(configString, namespace, ks); err != nil {
		t.Errorf("ApplyConfigFromString(%s, %s, ks) = %v; want <nil>", configString, namespace, err)
	}
}

func TestApplyConfigFromStringErrors(t *testing.T) {
	configString := string(fileContents(t, "testing/deployment.yaml"))
	namespace := "default"
	ks := &testservices.TestKubectl{
		ApplyFromStringResponse: map[string][]error{
			configString: {fmt.Errorf("failed to apply kubernetes manifests to cluster")},
		},
	}

	if err := ApplyConfigFromString(configString, namespace, ks); err == nil {
		t.Errorf("ApplyConfigs(%s, %s, ks) = <nil>; want error", configString, namespace)
	}
}

func TestGetDeployedObject(t *testing.T) {
	ctx := context.Background()

	testDeploymentFile := "testing/deployment.yaml"
	testServiceFile := "testing/service.yaml"

	tests := []struct {
		name string

		kind      string
		objName   string
		namespace string
		ks        services.KubectlService

		want runtime.Object
	}{{
		name: "Get deployed deployment",

		kind:      "Deployment",
		objName:   "test-app",
		namespace: "default",
		ks: &testservices.TestKubectl{
			GetResponse: map[string]map[string][]testservices.GetResponse{
				"Deployment": {
					"test-app": {
						{
							Res: string(fileContents(t, testDeploymentFile)),
							Err: nil,
						},
					},
				},
			},
		},

		want: newObjectFromFile(t, testDeploymentFile),
	}, {
		name: "Get deployed service",

		kind:      "Service",
		objName:   "test-app",
		namespace: "default",
		ks: &testservices.TestKubectl{
			GetResponse: map[string]map[string][]testservices.GetResponse{
				"Service": {
					"test-app": {
						{
							Res: string(fileContents(t, testServiceFile)),
							Err: nil,
						},
					},
				},
			},
		},

		want: newObjectFromFile(t, testServiceFile),
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := GetDeployedObject(ctx, tc.kind, tc.objName, tc.namespace, tc.ks); !reflect.DeepEqual(got, tc.want) || err != nil {
				t.Errorf("GetDeployedObject(ctx, %s, %s, %s, ks,) = %v, %v; want %v, <nil>", tc.kind, tc.objName, tc.namespace, got, err, tc.want)
			}
		})
	}
}

func TestDeployedObjectExists(t *testing.T) {
	ctx := context.Background()

	testDeploymentFile := "testing/deployment.yaml"

	tests := []struct {
		name string

		kind      string
		objName   string
		namespace string
		ks        services.KubectlService

		want bool
	}{{
		name: "Deployed deployment exists",

		kind:      "Deployment",
		objName:   "test-app",
		namespace: "default",
		ks: &testservices.TestKubectl{
			GetResponse: map[string]map[string][]testservices.GetResponse{
				"Deployment": {
					"test-app": {
						{
							Res: string(fileContents(t, testDeploymentFile)),
							Err: nil,
						},
					},
				},
			},
		},

		want: true,
	}, {
		name: "Deployed service does not exist",

		kind:      "Service",
		objName:   "test-app",
		namespace: "default",
		ks: &testservices.TestKubectl{
			GetResponse: map[string]map[string][]testservices.GetResponse{
				"Service": {
					"test-app": {
						{
							Res: "",
							Err: nil,
						},
					},
				},
			},
		},

		want: false,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := DeployedObjectExists(ctx, tc.kind, tc.objName, tc.namespace, tc.ks); got != tc.want || err != nil {
				t.Errorf("DeployedObjectExists(ctx, %s, %s, %s, ks,) = %t, %v; want %t, <nil>", tc.kind, tc.objName, tc.namespace, got, err, tc.want)
			}
		})
	}
}

func TestDeployedObjectExistsErrors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		kind      string
		objName   string
		namespace string
		ks        services.KubectlService
	}{{
		name: "Failed to check if object exists",

		kind:      "Service",
		objName:   "test-app",
		namespace: "default",
		ks: &testservices.TestKubectl{
			GetResponse: map[string]map[string][]testservices.GetResponse{
				"Service": {
					"test-app": {
						{
							Res: "",
							Err: fmt.Errorf("failed to get service"),
						},
					},
				},
			},
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := DeployedObjectExists(ctx, tc.kind, tc.objName, tc.namespace, tc.ks); err == nil {
				t.Errorf("DeployedObjectExists(ctx, %s, %s, %s, ks,) = %t, <nil>; want false, error", tc.kind, tc.objName, tc.namespace, got)
			}
		})
	}
}

func newObjectFromFile(t *testing.T, filename string) runtime.Object {
	contents := fileContents(t, filename)
	obj, err := resource.DecodeFromYAML(nil, contents)
	if err != nil {
		t.Fatalf("failed to decode resource from file %s", filename)
	}
	return obj
}

func fileContents(t *testing.T, filename string) []byte {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read file %s", filename)
	}
	return contents
}
