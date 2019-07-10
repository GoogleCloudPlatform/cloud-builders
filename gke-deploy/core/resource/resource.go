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
package resource

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/google/go-containerregistry/pkg/name"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/image"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
)

type resourceDecoder struct{}

func (resourceDecoder) Decode(data []byte, defaults *schema.GroupVersionKind, into runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	json, err := yaml.ToJSON(data)
	if err != nil {
		return nil, nil, err
	}
	return unstructured.UnstructuredJSONScheme.Decode(json, defaults, into)
}

var (
	decoder = resourceDecoder{}
	encoder = json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{Yaml: true})
)

// Objects maps resource file base names to corresponding resource objects (mutable).
type Objects map[string]*Object

// Object extends unstructured.Unstructured, which implements runtime.Object
type Object struct {
	*unstructured.Unstructured
}

// DecodeFromYAML decodes an object from a YAML string as bytes.
func DecodeFromYAML(ctx context.Context, yaml []byte) (*Object, error) {
	obj, err := runtime.Decode(decoder, yaml)
	if err != nil {
		return nil, fmt.Errorf("failed to decode yaml into object")
	}
	objUn, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("failed to convert object to Unstructured")
	}
	return &Object{
		objUn,
	}, nil
}

// ParseConfigs parses resource objects from a file or directory of files into a map that maps
// unique file base names to the parsed objects.
func ParseConfigs(ctx context.Context, configs string, oss services.OSService) (Objects, error) {
	objs := Objects{}

	fi, err := oss.Stat(ctx, configs)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %q: %v", configs, err)
	}

	if fi.IsDir() {
		d := configs
		files, err := oss.ReadDir(ctx, d)
		if err != nil {
			return nil, fmt.Errorf("failed to list files in directory %q: %v", d, err)
		}

		hasResource := false
		for _, fi := range files {
			if fi.IsDir() {
				continue
			}
			filename := filepath.Join(d, fi.Name())
			if !hasYamlOrYmlSuffix(filename) {
				continue
			}

			if err := parseResourcesFromFile(ctx, filename, objs, oss); err != nil {
				return nil, fmt.Errorf("failed to parse config %q: %v", filename, err)
			}
			hasResource = true
		}
		if !hasResource {
			return nil, fmt.Errorf("directory %q has no \".yaml\" or \".yaml\" files to parse", d)
		}
	} else {
		filename := configs
		if !hasYamlOrYmlSuffix(filename) {
			return nil, fmt.Errorf("file %q does not end in \".yaml\" or \".yml\"", filename)
		}

		if err := parseResourcesFromFile(ctx, filename, objs, oss); err != nil {
			return nil, fmt.Errorf("failed to parse config %q: %v", filename, err)
		}
	}

	return objs, nil
}

// SaveAsConfigs saves resource objects as config files to a target output directory.
func SaveAsConfigs(ctx context.Context, objs Objects, outputDir string, oss services.OSService) error {
	fi, err := oss.Stat(ctx, outputDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to get file info for output directory %q: %v", outputDir, err)
	}

	if err == nil && !fi.IsDir() {
		return fmt.Errorf("output directory %q exists as a file", outputDir)
	}

	if err == nil && fi.IsDir() {
		files, err := oss.ReadDir(ctx, outputDir)
		if err != nil {
			return fmt.Errorf("failed to list files in output directory %q: %v", outputDir, err)
		}
		if len(files) != 0 {
			return fmt.Errorf("output directory %q exists and is not empty", outputDir)
		}
	}

	if err := oss.MkdirAll(ctx, outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory %q: %v", outputDir, err)
	}
	for baseName, obj := range objs {
		filename := filepath.Join(outputDir, baseName)
		out, err := runtime.Encode(encoder, obj)
		if err != nil {
			return fmt.Errorf("failed to encode resource: %v", err)
		}
		if err := oss.WriteFile(ctx, filename, out, 0644); err != nil {
			return fmt.Errorf("failed to write file %q: %v", filename, err)
		}
	}
	return nil
}

// UpdateMatchingContainerImage updates all objects that have container images matching the provided image
// name with the provided replacement string.
func UpdateMatchingContainerImage(ctx context.Context, objs Objects, imageName, replace string) error {
	matched := false
	for _, obj := range objs {
		var nestedFields []string

		switch kind := ResourceKind(obj); kind {
		case "CronJob":
			nestedFields = []string{"spec", "jobTemplate", "spec", "template", "spec", "containers"}
		case "Pod":
			nestedFields = []string{"spec", "containers"}
		case "DaemonSet", "Deployment", "Job", "ReplicaSet", "ReplicationController", "StatefulSet":
			nestedFields = []string{"spec", "template", "spec", "containers"}
		default:
			continue
		}

		cons, ok, err := unstructured.NestedFieldNoCopy(obj.Object, nestedFields...)
		if err != nil {
			return fmt.Errorf("failed to get nested containers field: %v", err)
		}
		if !ok {
			continue
		}
		consList, ok := cons.([]interface{})
		if !ok {
			return fmt.Errorf("failed to convert containers to list")
		}

		for _, con := range consList {
			conMap, ok := con.(map[string]interface{})
			if !ok {
				return fmt.Errorf("failed to convert container to map")
			}
			im, ok, err := unstructured.NestedString(conMap, "image")
			if err != nil {
				return fmt.Errorf("failed to get image field: %v", err)
			}
			if !ok {
				continue
			}

			ref, err := name.ParseReference(im)
			if err != nil {
				return fmt.Errorf("failed to parse reference from image %q: %v", im, err)
			}
			if image.Name(ref) == imageName {
				fmt.Printf("Updating container of resource: %v\n", obj)
				if err := unstructured.SetNestedField(conMap, replace, "image"); err != nil {
					return fmt.Errorf("failed to set image field: %v", err)
				}
				matched = true
			}
		}
	}

	if !matched {
		fmt.Fprintf(os.Stderr, "\nWARNING: Did not find any resources with a container that has image name %q\n\n", imageName)
	}

	return nil
}

// UpdateNamespace updates all objects to change its namespace to the provided namespace. Objects
// that do not have a namespace field will not be updated.
func UpdateNamespace(ctx context.Context, objs Objects, replace string) error {
	var hasNS []*Object
	for _, obj := range objs {
		ns, err := resourceNamespace(obj)
		if err != nil {
			return fmt.Errorf("failed to get namespace field: %v", err)
		}
		if ns != "" {
			if err := setResourceNamespace(obj, replace); err != nil {
				return fmt.Errorf("failed to set namespace field: %v", err)
			}
			hasNS = append(hasNS, obj)
		}
	}
	if len(hasNS) > 0 {
		fmt.Fprintf(os.Stderr, "\nWARNING: It is recommended to set a resource's namespace at deploy time, rather than in its config. The following resources have embedded namespaces:\n")
		for _, obj := range hasNS {
			fmt.Fprintf(os.Stderr, "%v\n", obj)
		}
		fmt.Fprintln(os.Stderr)
	}
	return nil
}

// HasObject returns true if there exists an object in objs that matches the provided kind and
// name.
func HasObject(ctx context.Context, objs Objects, kind, name string) (bool, error) {
	for _, obj := range objs {
		objKind := ResourceKind(obj)
		objName, err := ResourceName(obj)
		if err != nil {
			return false, fmt.Errorf("failed to get resource name: %v", err)
		}
		if objKind == kind && objName == name {
			return true, nil
		}
	}
	return false, nil
}

// AddObject adds the provided object to objs with a generated file base name as its key.
func AddObject(ctx context.Context, objs Objects, obj *Object) error {
	// Try <resource-kind>.yaml
	objKind := strings.ToLower(ResourceKind(obj))
	baseName := fmt.Sprintf("%s.yaml", objKind)
	if _, ok := objs[baseName]; !ok {
		objs[baseName] = obj
		return nil
	}

	// Try <resource-kind>-<resource-name>.yaml
	objName, err := ResourceName(obj)
	if err != nil {
		return fmt.Errorf("failed to get resource name: %v", err)
	}
	baseName = fmt.Sprintf("%s-%s.yaml", objKind, objName)
	if _, ok := objs[baseName]; !ok {
		objs[baseName] = obj
		return nil
	}

	// Try <resource-kind>-<resource-name>-#.yaml
	fixedBaseName, err := fixCollidingFileBaseName(baseName, objs)
	if err != nil {
		return fmt.Errorf("failed to fix colliding base name %q: %v", baseName, err)
	}
	objs[fixedBaseName] = obj
	return nil
}

// CreateNamespaceObject creates a namespace object with the given name.
func CreateNamespaceObject(ctx context.Context, name string) (*Object, error) {
	if name == "default" {
		return nil, fmt.Errorf("namespace name should not be \"default\"")
	}
	obj, err := DecodeFromYAML(ctx, []byte(namespaceTemplate))
	if err != nil {
		return nil, fmt.Errorf("failed to create template namespace object")
	}

	if err := setResourceName(obj, name); err != nil {
		return nil, fmt.Errorf("failed to set name field: %v", err)
	}
	return obj, nil
}

// DeploySummary returns a string representation of a summary of a list of objects' deploy statuses.
func DeploySummary(ctx context.Context, objs Objects) (string, error) {
	// Sort values
	var sorted []*Object
	for _, obj := range objs {
		sorted = append(sorted, obj)
	}
	sorted = sortObjectsByKindAndName(sorted)

	// Create table
	padding := 4
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, padding, ' ', 0)

	if _, err := fmt.Fprintln(w, "KIND\tNAME\tREADY\t"); err != nil {
		return "", fmt.Errorf("failed to write to writer: %v", err)
	}

	for _, obj := range sorted {
		kind := ResourceKind(obj)
		name, err := ResourceName(obj)
		if err != nil {
			return "", fmt.Errorf("failed to get resource name: %v", err)
		}

		extraInfo, err := deploySummaryExtraInfo(obj)
		if err != nil {
			return "", fmt.Errorf("failed to get resource summary extra info: %v", err)
		}

		var ready string
		ok, err := IsReady(ctx, obj)
		if err != nil {
			ready = "Unknown"
		} else if ok {
			ready = "Yes"
		} else {
			ready = "No"
		}

		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", kind, name, ready, extraInfo); err != nil {
			return "", fmt.Errorf("failed to write to writer: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush writer: %v", err)
	}

	return buf.String(), nil
}

// sortObjectsByKindAndName sorts a list of objects by kind, then name, alphabetically.
func sortObjectsByKindAndName(objs []*Object) []*Object {
	sort.SliceStable(objs, func(i, j int) bool {
		a := objs[i]
		b := objs[j]

		aKind := ResourceKind(a)
		bKind := ResourceKind(b)
		aName, err := ResourceName(a)
		if err != nil {
			return false // Move a to end of slice
		}
		bName, err := ResourceName(b)
		if err != nil {
			return true // Move b to end of slice
		}

		if aKind == bKind {
			return aName < bName
		}
		return aKind < bKind
	})
	return objs
}

func deploySummaryExtraInfo(obj *Object) (string, error) {
	var extraInfo string

	kind := ResourceKind(obj)
	switch kind {
	case "Service":
		serviceType, ok, err := unstructured.NestedString(obj.Object, "spec", "type")
		if err != nil {
			return "", fmt.Errorf("failed to get spec.type field: %v", err)
		}
		if !ok || serviceType == "" {
			return "", fmt.Errorf("spec.type field is missing or is empty")
		}
		switch serviceType {
		case "LoadBalancer":
			return serviceIPs(obj)
		case "ExternalName":
			return serviceExternalName(obj)
		}
	default:
	}

	return extraInfo, nil
}

func serviceIPs(obj *Object) (string, error) {
	ingress, ok, err := unstructured.NestedSlice(obj.Object, "status", "loadBalancer", "ingress")
	if err != nil {
		return "", fmt.Errorf("failed to get status.loadBalancer.ingress field: %v", err)
	}
	if !ok || len(ingress) == 0 {
		return "", nil
	}

	var ips []string
	for _, i := range ingress {
		iMap, ok := i.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("failed to convert ingress to map")
		}
		ip, ok, err := unstructured.NestedString(iMap, "ip")
		if err != nil {
			return "", fmt.Errorf("failed to get ip field: %v", err)
		}
		if !ok || ip == "" {
			return "", fmt.Errorf("ip field is missing or is empty")
		}
		ips = append(ips, ip)
	}

	return strings.Join(ips, ", "), nil
}

func serviceExternalName(obj *Object) (string, error) {
	externalName, ok, err := unstructured.NestedString(obj.Object, "spec", "externalName")
	if err != nil {
		return "", fmt.Errorf("failed to get spec.externalName field: %v", err)
	}
	if !ok {
		return "", fmt.Errorf("spec.externalName field is missing")
	}

	return externalName, nil
}

func parseResourcesFromFile(ctx context.Context, filename string, objs Objects, oss services.OSService) error {
	in, err := oss.ReadFile(ctx, filename)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %v", filename, err)
	}

	split := strings.Split(string(in), "\n---")

	for i, r := range split {
		obj, err := DecodeFromYAML(ctx, []byte(r))
		if err != nil {
			return fmt.Errorf("failed to decode resource from item %d in file %q: %v", i+1, filename, err)
		}

		// For configs containing one resource, just use the original file base name.
		// e.g., if ".../deployment.yaml" contains one resource, the resource should be given the file base name "deployment.yaml".
		baseName := filepath.Base(filename)

		// For configs containing multiple resources, each resource is given the file base  name <file-prefix>-<resource-kind>-<resource-name>.<file-suffix>.
		// e.g., if ".../resource.yaml" contains multiple resources, with one being a Deployment with the name "nginx",
		// the resource will be given the file base name "resource-deployment-nginx.yaml".
		if len(split) > 1 {
			// This is the case where the file has more than one resource, separated by "---".
			ix := strings.LastIndex(baseName, ".")
			prefix := baseName[:ix]
			suffix := baseName[ix+1:]
			objKind := strings.ToLower(ResourceKind(obj))
			objName, err := ResourceName(obj)
			if err != nil {
				return fmt.Errorf("failed to get resource name of item %d in file %q: %v", i+1, filename, err)
			}
			baseName = fmt.Sprintf("%s-%s-%s.%s", prefix, objKind, objName, suffix)
		}

		fixedBaseName, err := fixCollidingFileBaseName(baseName, objs)
		if err != nil {
			return fmt.Errorf("failed to fix colliding base name %q: %v", baseName, err)
		}

		objs[fixedBaseName] = obj
	}

	return nil
}

func fixCollidingFileBaseName(name string, objs Objects) (string, error) {
	if _, ok := objs[name]; !ok {
		return name, nil
	}

	const max = 1000
	var newName string
	for i := 2; i < max; i++ {
		x := strings.LastIndex(name, ".")
		prefix := name[:x]
		suffix := name[x+1:]
		newName = fmt.Sprintf("%s-%d.%s", prefix, i, suffix)
		if _, ok := objs[newName]; !ok {
			return newName, nil
		}
	}
	return "", fmt.Errorf("reached upper limit %d", max)
}

// String returns a string representation of objects.
func (objs Objects) String() string {
	// Sort values
	var sorted []*Object
	for _, obj := range objs {
		sorted = append(sorted, obj)
	}
	sorted = sortObjectsByKindAndName(sorted)
	return fmt.Sprintf("%v", sorted)
}

// String returns a string representation of an object.
func (obj *Object) String() string {
	kind := ResourceKind(obj)
	name, err := ResourceName(obj)
	if err != nil {
		name = "UNKNOWN"
	}

	return fmt.Sprintf("{kind: %s, name: %s}", kind, name)
}

// AddLabel updates an object to add a label with the key and value provided.
func AddLabel(ctx context.Context, obj *Object, key, value string, override bool) error {
	if key == "" || value == "" {
		return fmt.Errorf("key and value cannot be empty")
	}

	if err := addLabelToNestedField(obj, key, value, override, "metadata", "labels"); err != nil {
		return err
	}

	var nestedFields []string
	switch kind := ResourceKind(obj); kind {
	case "CronJob":
		nestedFields = []string{"spec", "jobTemplate", "spec", "template", "metadata", "labels"}
	case "DaemonSet", "Deployment", "Job", "ReplicaSet", "ReplicationController", "StatefulSet":
		nestedFields = []string{"spec", "template", "metadata", "labels"}
	default:
		return nil
	}
	if err := addLabelToNestedField(obj, key, value, override, nestedFields...); err != nil {
		return err
	}

	return nil
}

func addLabelToNestedField(obj *Object, key, value string, override bool, nestedFields ...string) error {
	labels, ok, err := unstructured.NestedMap(obj.Object, nestedFields...)
	if err != nil {
		return fmt.Errorf("failed to get labels field: %v", err)
	}

	if !ok {
		labels = make(map[string]interface{})
	}

	if existing, ok := labels[key]; ok && !override {
		if existing != value {
			fmt.Fprintf(os.Stderr, "\nWARNING: Label %q is already set as %q for object %v in %v field. Not overriding.\n", key, existing, obj, strings.Join(nestedFields, "."))
		}
	} else {
		labels[key] = value
		if err := unstructured.SetNestedMap(obj.Object, labels, nestedFields...); err != nil {
			return fmt.Errorf("failed to set labels field: %v", err)
		}
	}

	return nil
}

// TODO(joonlim): These should be member functions of Object.

func ResourceKind(obj *Object) string {
	return obj.GetObjectKind().GroupVersionKind().Kind
}

func ResourceName(obj *Object) (string, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return "", fmt.Errorf("failed to get metadata accessor from object: %v", err)
	}
	return accessor.GetName(), nil
}

func setResourceName(obj *Object, name string) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return fmt.Errorf("failed to get metadata accessor from object: %v", err)
	}
	accessor.SetName(name)
	return nil
}

func resourceNamespace(obj *Object) (string, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return "", fmt.Errorf("failed to get metadata accessor from object: %v", err)
	}
	return accessor.GetNamespace(), nil
}

func setResourceNamespace(obj *Object, namespace string) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return fmt.Errorf("failed to get metadata accessor from object: %v", err)
	}
	accessor.SetNamespace(namespace)
	return nil
}

func hasYamlOrYmlSuffix(filename string) bool {
	return strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml")
}
