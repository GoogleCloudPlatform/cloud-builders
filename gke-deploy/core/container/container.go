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
package container

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
	"github.com/google/go-containerregistry/pkg/name"
)

// ParseImages parses a slice of image strings.
func ParseImages(images []string) ([]name.Reference, error) {
	var refs []name.Reference

	exists := make(map[string]bool)
	for _, image := range images {
		ref, err := parseImage(image)
		if err != nil {
			return nil, err
		}

		imName, err := GetName(ref)
		if err != nil {
			return nil, fmt.Errorf("failed to get image name: %v", err)
		}

		if ok := exists[imName]; ok {
			return nil, fmt.Errorf("duplicate image name: %q", imName)
		}
		exists[imName] = true
		refs = append(refs, ref)
	}

	return refs, nil
}

func parseImage(image string) (name.Reference, error) {
	im, err := name.ParseReference(image)
	if err != nil {
		return nil, fmt.Errorf("image is invalid: %q", image)
	}
	return im, nil
}

// ParseImageReference gets an image's name, given a string representation of the image.
// e.g., given "gcr.io/my-project/my-app:1.0.0", returns "gcr.io/project/my-app"
func ParseImageReference(image string) (string, error) {
	im, err := parseImage(image)
	if err != nil {
		return "", err
	}
	return GetName(im)
}

// GetName gets an image's name.
func GetName(image name.Reference) (string, error) {
	switch t := image.(type) {
	case name.Tag:
		return fmt.Sprintf("%s/%s", t.RegistryStr(), t.RepositoryStr()), nil
	case name.Digest:
		return fmt.Sprintf("%s/%s", t.RegistryStr(), t.RepositoryStr()), nil
	default:
		return "", fmt.Errorf("invalid image type: %s", t)
	}
}

// GetDigest gets an image's corresponding digest.
func GetDigest(ctx context.Context, image name.Reference, rs services.RemoteService) (string, error) {
	im, err := rs.Image(image)
	if err != nil {
		return "", fmt.Errorf("failed to get remote image reference: %v", err)
	}
	digest, err := im.Digest()
	if err != nil {
		return "", fmt.Errorf("failed to get image digest: %v", err)
	}
	return fmt.Sprintf("%s:%s", digest.Algorithm, digest.Hex), nil
}
