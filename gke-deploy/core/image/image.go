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
// Package image contains related to container images.
package image

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
	"github.com/google/go-containerregistry/pkg/name"
)

// Name gets an image's name from a Reference.
// e.g., If the string representation of the Reference is "gcr.io/my-project/my-image:1.0.0", this
// returns "gcr.io/my-project/my-image".
func Name(ref name.Reference) string {
	return fmt.Sprintf("%s/%s", ref.Context().RegistryStr(), ref.Context().RepositoryStr())
}

// ResolveDigest gets an image's corresponding digest.
func ResolveDigest(ctx context.Context, ref name.Reference, rs services.RemoteService) (string, error) {
	im, err := rs.Image(ref)
	if err != nil {
		return "", fmt.Errorf("failed to get remote image reference: %v", err)
	}
	digest, err := im.Digest()
	if err != nil {
		return "", fmt.Errorf("failed to get image digest: %v", err)
	}
	return fmt.Sprintf("%s:%s", digest.Algorithm, digest.Hex), nil
}
