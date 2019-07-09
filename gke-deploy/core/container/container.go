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
)

const (
	imageDigestFormat = "value(image_summary.digest)"
)

// GetDigest uses an image string to get its corresponding digest.
func GetDigest(ctx context.Context, image string, gs services.GcloudService) (string, error) {
	digest, err := gs.ContainerImagesDescribe(ctx, image, imageDigestFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get image digest: %v", err)
	}
	return digest, nil
}
