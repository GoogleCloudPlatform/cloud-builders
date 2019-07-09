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

const (
	imageDigestFormat = "value(image_summary.digest)"
)

// Image is an alias to name.Reference.
type Image = name.Reference

// ParseImages parses a slice of image strings.
func ParseImages(ctx context.Context, images []string) ([]Image, error) {
	var ims []Image

	exists := make(map[string]bool)
	for _, image := range images {
		im, err := parseImage(ctx, image)
		if err != nil {
			return nil, err
		}

		imName, err := GetName(ctx, im)
		if err != nil {
			return nil, fmt.Errorf("failed to get image name: %v", err)
		}

		if ok := exists[imName]; ok {
			return nil, fmt.Errorf("duplicate image name: %q", imName)
		}
		exists[imName] = true
		ims = append(ims, im)
	}

	return ims, nil
}

func parseImage(ctx context.Context, image string) (Image, error) {
	im, err := name.ParseReference(image)
	if err != nil {
		return nil, fmt.Errorf("image is invalid: %q", image)
	}
	return im, nil
}

// GetNameFromString gets an image's name, given a string representation of the image.
func GetNameFromString(ctx context.Context, image string) (string, error) {
	im, err := parseImage(ctx, image)
	if err != nil {
		return "", err
	}
	return GetName(ctx, im)
}

// GetName gets an image's name.
func GetName(ctx context.Context, image Image) (string, error) {
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
func GetDigest(ctx context.Context, image Image, gs services.GcloudService) (string, error) {
	switch t := image.(type) {
	case name.Tag:
		digest, err := gs.ContainerImagesDescribe(ctx, t.Name(), imageDigestFormat)
		if err != nil {
			return "", fmt.Errorf("failed to get image digest: %v", err)
		}
		return digest, nil
	case name.Digest:
		return image.Identifier(), nil
	default:
		return "", fmt.Errorf("invalid image type: %s", t)
	}
}
