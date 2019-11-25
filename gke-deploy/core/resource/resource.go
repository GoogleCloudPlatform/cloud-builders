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
// Package resource contains logic related to Kubernetes resource objects.
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

// AggregatedFilename is the filename for the file created by SaveAsConfigs.
const AggregatedFilename = "aggregated-resources.yaml"

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
type Objects []*Object

// Object extends unstructured.Unstructured, which implements runtime.Object
type Object struct {
	*unstructured.Unstructured
}

// EncodeToYAMLString encodes an object from *Object to a string.
func EncodeToYAMLString(obj *Object) (string, error) {
	out, err := runtime.Encode(encoder, obj)
	if err != nil {
		return "", fmt.Errorf("failed to encode resource: %v", err)
	}
	return string(out), nil
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
func ParseConfigs(ctx context.Context, configs string, oss services.OSService, recursive bool) (Objects, error) {
	objs := Objects{}

	if configs == "-" {
		if recursive {
			return nil, fmt.Errorf("cannot recur with stdin")
		}
		objs, err := parseResourcesFromFile(ctx, configs, objs, oss)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config from stdin: %v", err)
		}
		return objs, nil
	}

	fi, err := oss.Stat(ctx, configs)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %q: %v", configs, err)
	}

	if !fi.IsDir() && recursive {
		return nil, fmt.Errorf("cannot recur through a file")
	}

	hasResources := false

	// Since walk is recursive, we need to declare it before creating the function that refers to it.
	var walk func(path string, fi os.FileInfo, baseDir bool) error
	walk = func(path string, fi os.FileInfo, baseDir bool) error {
		if fi.IsDir() {
			if !baseDir && !recursive {
				return nil
			}
			subfiles, err := oss.ReadDir(ctx, path)
			if err != nil {
				return fmt.Errorf("failed to list files in directory %q: %v", path, err)
			}
			for _, subfile := range subfiles {
				subpath := filepath.Join(path, subfile.Name())
				if err = walk(subpath, subfile, false); err != nil {
					return err
				}
			}
		} else {
			if hasYamlOrYmlSuffix(path) {
				hasResources = true
				objs, err = parseResourcesFromFile(ctx, path, objs, oss)
				if err != nil {
					return fmt.Errorf("failed to parse config %q: %v", path, err)
				}
			}
		}
		return nil
	}

	if err = walk(configs, fi, true); err != nil {
		return nil, err
	}

	if !hasResources {
		if fi.IsDir() {
			return nil, fmt.Errorf("directory %q has no \".yaml\" or \".yml\" files to parse", configs)
		}
		return nil, fmt.Errorf("file %q does not end in \".yaml\" or \".yml\"", configs)
	}

	return objs, nil
}

// SaveAsConfigs saves resource objects as config files to a target output directory.
// If any lines in a resource object's string representation contain a key in
// lineComments, the corresponding value will be added as a comment at the end of
// the line.
func SaveAsConfigs(ctx context.Context, objs Objects, outputDir string, lineComments map[string]string, oss services.OSService) error {
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

	filename := filepath.Join(outputDir, AggregatedFilename)

	resources := make([]string, 0, len(objs))

	for _, obj := range objs {
		out, err := runtime.Encode(encoder, obj)
		if err != nil {
			return fmt.Errorf("failed to encode resource: %v", err)
		}

		outWithComments, err := addCommentsToLines(string(out), lineComments)
		if err != nil {
			return fmt.Errorf("failed to add comment to object file: %v", err)
		}

		resources = append(resources, outWithComments)
	}

	contents := strings.Join(resources, "\n\n---\n\n")
	if err := oss.WriteFile(ctx, filename, []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to write file %q: %v", filename, err)
	}
	return nil
}

// addCommentsToLines iterates through the lines of a string ('-n'-delimited)
// and if any lines contain a key in lineComments, the corresponding value will
// be added as a comment at the end of the line. This function returns the
// modified string.
func addCommentsToLines(s string, lineComments map[string]string) (string, error) {
	lines := strings.Split(s, "\n")
	lineIdx := 0
	for _, line := range lines {
		for stringToContain, comment := range lineComments {
			if strings.Contains(stringToContain, "\n") {
				return "", fmt.Errorf("line cannot contain a newline character")
			}
			if strings.Contains(comment, "\n") {
				return "", fmt.Errorf("comment cannot contain a newline character")
			}
			if strings.Contains(line, stringToContain) {
				lines[lineIdx] = fmt.Sprintf("%s  # %s", line, comment)
			}
		}
		lineIdx++
	}

	return strings.Join(lines, "\n"), nil
}

// UpdateMatchingContainerImage updates all objects that have container images matching the provided image
// name with the provided replacement string.
func UpdateMatchingContainerImage(ctx context.Context, objs Objects, imageName, replace string) error {
	matched := false
	for _, obj := range objs {
		var nestedFields []string

		switch kind := ObjectKind(obj); kind {
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
// that do not have a namespace field will also be updated to have a namespace field.
func UpdateNamespace(ctx context.Context, objs Objects, replace string) error {
	for _, obj := range objs {
		if err := setObjectNamespace(obj, replace); err != nil {
			return fmt.Errorf("failed to set namespace field: %v", err)
		}
	}
	return nil
}

// AddNamespaceIfMissing updates all objects to add a namespace only if the object does
// not have one already.
func AddNamespaceIfMissing(objs Objects, namespace string) error {
	for _, obj := range objs {
		ns, err := ObjectNamespace(obj)
		if err != nil {
			return fmt.Errorf("failed to get namespace field: %v", err)
		}
		if ns != "" {
			continue
		}
		if err := setObjectNamespace(obj, namespace); err != nil {
			return fmt.Errorf("failed to set namespace field: %v", err)
		}
	}
	return nil
}

// HasObject returns true if there exists an object in objs that matches the provided kind and
// name.
func HasObject(ctx context.Context, objs Objects, kind, name string) (bool, error) {
	for _, obj := range objs {
		objKind := ObjectKind(obj)
		objName, err := ObjectName(obj)
		if err != nil {
			return false, fmt.Errorf("failed to get resource name: %v", err)
		}
		if objKind == kind && objName == name {
			return true, nil
		}
	}
	return false, nil
}

// CreateDeploymentObject creates a Deployment object with the given name and image.
// The created Deployment will have 3 replicas.
func CreateDeploymentObject(ctx context.Context, name string, selectorValue, image string) (*Object, error) {
	obj, err := DecodeFromYAML(ctx, []byte(fmt.Sprintf(deploymentTemplate, name, "app", selectorValue, "app", selectorValue, name, image)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode Deployment object from template")
	}
	return obj, nil
}

// CreateHorizontalPodAutoscalerObject creates a Namespace object with the given name.
// The created HorizontalPodAutoscaler will have minReplicas set to 1, maxReplicas set to 5, and a
// cpu targetAverageUtilization of 80.
func CreateHorizontalPodAutoscalerObject(ctx context.Context, name, deploymentName string) (*Object, error) {
	obj, err := DecodeFromYAML(ctx, []byte(fmt.Sprintf(horizontalPodAutoscalerTemplate, name, deploymentName)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode HorizontalPodAutoscaler object from template")
	}
	return obj, nil
}

// CreateNamespaceObject creates a Namespace object with the given name.
func CreateNamespaceObject(ctx context.Context, name string) (*Object, error) {
	if name == "default" {
		return nil, fmt.Errorf("namespace name should not be \"default\"")
	}
	obj, err := DecodeFromYAML(ctx, []byte(fmt.Sprintf(namespaceTemplate, name)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode Namespace object from template")
	}
	return obj, nil
}

// CreateServiceObject creates a Service object with the given name with service type LoadBalancer.
func CreateServiceObject(ctx context.Context, name, selectorKey, selectorValue string, port int) (*Object, error) {
	obj, err := DecodeFromYAML(ctx, []byte(fmt.Sprintf(serviceTemplate, name, selectorKey, selectorValue, port, port)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode Service object from template")
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

	if _, err := fmt.Fprintln(w, "NAMESPACE\tKIND\tNAME\tREADY\t"); err != nil {
		return "", fmt.Errorf("failed to write to writer: %v", err)
	}

	for _, obj := range sorted {
		kind := ObjectKind(obj)
		name, err := ObjectName(obj)
		if err != nil {
			return "", fmt.Errorf("failed to get resource name: %v", err)
		}
		namespace, err := ObjectNamespace(obj)
		if err != nil {
			return "", fmt.Errorf("failed to get namespace of object: %v", err)
		}
		if namespace == "" {
			namespace = "default"
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

		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", namespace, kind, name, ready, extraInfo); err != nil {
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

		aKind := ObjectKind(a)
		bKind := ObjectKind(b)
		aName, err := ObjectName(a)
		if err != nil {
			return false // Move a to end of slice
		}
		bName, err := ObjectName(b)
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

	kind := ObjectKind(obj)
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
	ports, ok, err := unstructured.NestedSlice(obj.Object, "spec", "ports")
	if err != nil {
		return "", fmt.Errorf("failed to get spec.ports field: %v", err)
	}
	if !ok || len(ports) == 0 {
		return "", nil
	}
	portMap := ports[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to convert port to map")
	}
	port, ok, err := unstructured.NestedInt64(portMap, "port")
	if err != nil {
		return "", fmt.Errorf("failed to get port field: %v", err)
	}
	if !ok {
		return "", fmt.Errorf("port field is missing")
	}

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

		if port != 80 {
			ip = fmt.Sprintf("%s:%d", ip, port)
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

func parseResourcesFromFile(ctx context.Context, filename string, objs Objects, oss services.OSService) (Objects, error) {
	readStdin := filename == "-"
	var printFilename string
	if readStdin {
		printFilename = "stdin"
	} else {
		printFilename = fmt.Sprintf("file %q", filename)
	}

	in, err := oss.ReadFile(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", printFilename, err)
	}
	if readStdin {
		filename = "k8s.yaml" // Files parsed from stdin will have the prefix "k8s".
	}

	split := strings.Split(string(in), "\n---")

	for i, r := range split {
		lines := strings.Split(r, "\n")
		onlyCommentsAndWhitespace := true
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
				onlyCommentsAndWhitespace = false
				break
			}
		}
		if onlyCommentsAndWhitespace {
			continue
		}

		obj, err := DecodeFromYAML(ctx, []byte(r))
		if err != nil {
			return nil, fmt.Errorf("failed to decode resource from item %d in %s: %v", i+1, printFilename, err)
		}

		objs = append(objs, obj)
	}

	return objs, nil
}

// String returns a string representation of objects.
func (objs Objects) String() string {
	return fmt.Sprintf("%v", sortObjectsByKindAndName(objs))
}

// String returns a string representation of an object.
func (obj *Object) String() string {
	kind := ObjectKind(obj)
	name, err := ObjectName(obj)
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

	if err := addToNestedMap(obj, key, value, override, "metadata", "labels"); err != nil {
		return err
	}

	var nestedFields []string
	switch kind := ObjectKind(obj); kind {
	case "CronJob":
		nestedFields = []string{"spec", "jobTemplate", "spec", "template", "metadata", "labels"}
	case "DaemonSet", "Deployment", "Job", "ReplicaSet", "ReplicationController", "StatefulSet":
		nestedFields = []string{"spec", "template", "metadata", "labels"}
	default:
		return nil
	}
	if err := addToNestedMap(obj, key, value, override, nestedFields...); err != nil {
		return err
	}

	return nil
}

// AddAnnotation updates an object to add an annotation with the key and value
// provided.
func AddAnnotation(obj *Object, key, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("key and value cannot be empty")
	}

	if err := addToNestedMap(obj, key, value, true, "metadata", "annotations"); err != nil {
		return err
	}

	var nestedFields []string
	switch kind := ObjectKind(obj); kind {
	case "CronJob":
		nestedFields = []string{"spec", "jobTemplate", "spec", "template", "metadata", "annotations"}
	case "DaemonSet", "Deployment", "Job", "ReplicaSet", "ReplicationController", "StatefulSet":
		nestedFields = []string{"spec", "template", "metadata", "annotations"}
	default:
		return nil
	}
	if err := addToNestedMap(obj, key, value, true, nestedFields...); err != nil {
		return err
	}

	return nil
}

func addToNestedMap(obj *Object, key, value string, override bool, nestedFields ...string) error {
	mapField, ok, err := unstructured.NestedMap(obj.Object, nestedFields...)
	if err != nil {
		return fmt.Errorf("failed to get map field: %v", err)
	}

	if !ok {
		mapField = make(map[string]interface{})
	}

	if existing, ok := mapField[key]; ok && !override {
		if existing != value {
			fmt.Fprintf(os.Stderr, "\nWARNING: Key %q is already set as %q for object %v in %v field. Not overriding.\n", key, existing, obj, strings.Join(nestedFields, "."))
		}
	} else {
		mapField[key] = value
		if err := unstructured.SetNestedMap(obj.Object, mapField, nestedFields...); err != nil {
			return fmt.Errorf("failed to set map field: %v", err)
		}
	}

	return nil
}

// TODO(joonlim): These should be member functions of Object.

// ObjectKind returns the kind of an object.
func ObjectKind(obj *Object) string {
	return obj.GetObjectKind().GroupVersionKind().Kind
}

// ObjectName returns the name of an object.
func ObjectName(obj *Object) (string, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return "", fmt.Errorf("failed to get metadata accessor from object: %v", err)
	}
	return accessor.GetName(), nil
}

// ObjectNamespace returns the namespace of an object.
func ObjectNamespace(obj *Object) (string, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return "", fmt.Errorf("failed to get metadata accessor from object: %v", err)
	}
	return accessor.GetNamespace(), nil
}

func setObjectNamespace(obj *Object, namespace string) error {
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
