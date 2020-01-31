package resource

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	applicationsv1beta1 "github.com/kubernetes-sigs/application/pkg/apis/app/v1beta1"
)

func TestCreateApplicationObject(t *testing.T) {
	testApplicationFile := "testing/application.yaml"
	testApplicationNoVersionFile := "testing/application-no-version.yaml"
	testDeploymentFile := "testing/deployment.yaml"
	testServiceFile := "testing/service.yaml"
	testNamespaceFile := "testing/namespace.yaml"

	tests := []struct {
		name string

		applicationName   string
		selectorKey       string
		selectorValue     string
		descriptorType    string
		descriptorVersion string
		componentObjs     Objects

		want *Object
	}{{
		name: "Application with Deployment and Service",

		applicationName:   "test-name",
		selectorKey:       "foo",
		selectorValue:     "bar",
		descriptorType:    "test-type",
		descriptorVersion: "test-version",
		componentObjs: Objects{
			newObjectFromFile(t, testDeploymentFile),
			newObjectFromFile(t, testServiceFile),
		},

		want: newObjectFromFile(t, testApplicationFile),
	}, {
		name: "Application with Deployment and Service, no version",

		applicationName: "test-name",
		selectorKey:     "foo",
		selectorValue:   "bar",
		descriptorType:  "test-type",
		componentObjs: Objects{
			newObjectFromFile(t, testDeploymentFile),
			newObjectFromFile(t, testServiceFile),
		},

		want: newObjectFromFile(t, testApplicationNoVersionFile),
	}, {
		name: "Application with Deployment and Service, repeated",

		applicationName:   "test-name",
		selectorKey:       "foo",
		selectorValue:     "bar",
		descriptorType:    "test-type",
		descriptorVersion: "test-version",
		componentObjs: Objects{
			newObjectFromFile(t, testDeploymentFile),
			newObjectFromFile(t, testServiceFile),
			newObjectFromFile(t, testDeploymentFile),
			newObjectFromFile(t, testServiceFile),
		},

		want: newObjectFromFile(t, testApplicationFile),
	}, {
		name: "Application with Deployment and Service, Namespace ignored",

		applicationName:   "test-name",
		selectorKey:       "foo",
		selectorValue:     "bar",
		descriptorType:    "test-type",
		descriptorVersion: "test-version",
		componentObjs: Objects{
			newObjectFromFile(t, testDeploymentFile),
			newObjectFromFile(t, testServiceFile),
			newObjectFromFile(t, testNamespaceFile),
		},

		want: newObjectFromFile(t, testApplicationFile),
	}, {
		name: "Application with Deployment and Service, Application ignored",

		applicationName:   "test-name",
		selectorKey:       "foo",
		selectorValue:     "bar",
		descriptorType:    "test-type",
		descriptorVersion: "test-version",
		componentObjs: Objects{
			newObjectFromFile(t, testDeploymentFile),
			newObjectFromFile(t, testServiceFile),
			newObjectFromFile(t, testApplicationFile),
		},

		want: newObjectFromFile(t, testApplicationFile),
	}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// TODO(joonlim): Convert unit tests that use reflect.DeepEqual to use cmp.Equal/cmp.Diff style.
			if got, err := CreateApplicationObject(tc.applicationName, tc.selectorKey, tc.selectorValue, tc.descriptorType, tc.descriptorVersion, tc.componentObjs); err != nil {
				t.Errorf("CreateApplicationObject(%s, %s, %s, %s, %s, %v) returned error:\n%v", tc.applicationName, tc.selectorKey, tc.selectorValue, tc.descriptorType, tc.descriptorVersion, tc.componentObjs, err)
			} else if !cmp.Equal(got, tc.want) {
				t.Errorf("CreateApplicationObject(%s, %s, %s, %s, %s, %v) produced diff (-want +got):\n%s", tc.applicationName, tc.selectorKey, tc.selectorValue, tc.descriptorType, tc.descriptorVersion, tc.componentObjs, cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestSetApplicationLinks(t *testing.T) {
	testApplicationFile := "testing/application.yaml"
	testApplicationWithLinksFile := "testing/application-with-links.yaml"
	testApplicationWithLinks2File := "testing/application-with-links-2.yaml"

	tests := []struct {
		name string

		obj   *Object
		links []applicationsv1beta1.Link

		beforeUpdate *Object
		want         *Object
	}{{
		name: "Add links",

		obj: newObjectFromFile(t, testApplicationFile),
		links: []applicationsv1beta1.Link{
			{
				Description: "My Description",
				URL:         "https://my-link.com",
			},
			{
				Description: "1234",
				URL:         "5678",
			},
		},

		beforeUpdate: newObjectFromFile(t, testApplicationFile),
		want:         newObjectFromFile(t, testApplicationWithLinksFile),
	}, {
		name: "Add duplicate links",

		obj: newObjectFromFile(t, testApplicationFile),
		links: []applicationsv1beta1.Link{
			{
				Description: "My Description",
				URL:         "https://my-link.com",
			},
			{
				Description: "My Description",
				URL:         "https://my-link.com",
			},
		},

		beforeUpdate: newObjectFromFile(t, testApplicationFile),
		want:         newObjectFromFile(t, testApplicationWithLinks2File),
	}, {
		name: "Empty list",

		obj:   newObjectFromFile(t, testApplicationFile),
		links: nil,

		beforeUpdate: newObjectFromFile(t, testApplicationFile),
		want:         newObjectFromFile(t, testApplicationFile),
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := SetApplicationLinks(tc.obj, tc.links); err != nil {
				t.Errorf("SetApplicationLinks(%v, %v) returned error:\n%v", tc.beforeUpdate, tc.links, err)
			} else if !cmp.Equal(tc.obj, tc.want) {
				t.Errorf("SetApplicationLinks(%v, %v) produced diff (-want +got):\n%s", tc.beforeUpdate, tc.links, cmp.Diff(tc.want, tc.obj))
			}
		})
	}
}

func TestSetApplicationLinksErrors(t *testing.T) {
	testDeploymentFile := "testing/deployment.yaml"

	obj := newObjectFromFile(t, testDeploymentFile)
	links := []applicationsv1beta1.Link{
		{
			Description: "My Description",
			URL:         "https://my-link.com",
		},
		{
			Description: "1234",
			URL:         "5678",
		},
	}
	beforeUpdate := newObjectFromFile(t, testDeploymentFile)

	if err := SetApplicationLinks(obj, links); err == nil {
		t.Errorf("SetApplicationLinks(%v, %v) = %v; want error", beforeUpdate, links, err)
	}
}
