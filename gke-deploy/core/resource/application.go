package resource

import (
	"fmt"
	"sort"
	"strings"

	applicationsv1beta1 "github.com/kubernetes-sigs/application/pkg/apis/app/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// CreateApplicationObject creates an Application CR object with the given name and fields.
// appVersion may be empty, in which case the resulting Application CR will not have a value for spec.descriptor.version.
// To add links, to the Application, use the SetApplicationLinks function on the output object,
func CreateApplicationObject(name, selectorKey, selectorValue, descriptorType, descriptorVersion string, componentObjs Objects) (*Object, error) {
	componentKindsRemoveDups := make(map[metav1.GroupKind]bool)
	for _, obj := range componentObjs {
		kind := ObjectKind(obj)
		if kind == "Namespace" || kind == "Application" {
			continue
		}
		apiVersion := obj.GetAPIVersion() // e.g., v1, apps/v1, autoscaling/v2beta1
		apiVersionSplit := strings.Split(apiVersion, "/")
		var group string
		if len(apiVersionSplit) == 1 {
			group = "core"
		} else {
			group = apiVersionSplit[0]
		}
		componentKindsRemoveDups[metav1.GroupKind{
			Group: group,
			Kind:  kind,
		}] = true
	}
	componentKinds := make([]metav1.GroupKind, 0, len(componentKindsRemoveDups))
	for k := range componentKindsRemoveDups {
		componentKinds = append(componentKinds, k)
	}
	// Sort to make spec.componentKinds deterministic
	sort.SliceStable(componentKinds, func(i, j int) bool {
		a := componentKinds[i]
		b := componentKinds[j]
		if a.Group == b.Group {
			return a.Kind < b.Kind
		}
		return a.Group < b.Group
	})

	app := &applicationsv1beta1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "app.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: applicationsv1beta1.ApplicationSpec{
			ComponentGroupKinds: componentKinds,
			Descriptor: applicationsv1beta1.Descriptor{
				Type:    descriptorType,
				Version: descriptorVersion,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					selectorKey: selectorValue,
				},
			},
		},
	}

	asMap, err := convertApplicationToMap(app)
	if err != nil {
		return nil, err
	}

	return &Object{
		&unstructured.Unstructured{
			Object: asMap,
		},
	}, nil
}

// SetApplicationLinks sets a list of links to an Application object's spec.descriptor.links field.
func SetApplicationLinks(obj *Object, links []applicationsv1beta1.Link) error {
	if ObjectKind(obj) != "Application" {
		return fmt.Errorf("object must be an Application to add links")
	}
	app := &applicationsv1beta1.Application{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Unstructured.Object, app)
	if err != nil {
		return fmt.Errorf("failed to convert from unstructured")
	}
	app.Spec.Descriptor.Links = links

	asMap, err := convertApplicationToMap(app)
	if err != nil {
		return err
	}

	obj.Unstructured.Object = asMap
	return nil
}

func convertApplicationToMap(app *applicationsv1beta1.Application) (map[string]interface{}, error) {
	asMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(app)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to unstructured")
	}

	// The resulting map will have null/empty values for metadata.creationTimestamp and status.
	// We remove these from the map to make sure they are not included in the output YAML.
	metadata, ok := asMap["metadata"]
	if ok {
		metadataMap, ok := metadata.(map[string]interface{})
		if ok {
			_, ok := metadataMap["creationTimestamp"]
			if ok {
				delete(metadataMap, "creationTimestamp")
			}
		}
	}
	_, ok = asMap["status"]
	if ok {
		delete(asMap, "status")
	}
	return asMap, nil
}
