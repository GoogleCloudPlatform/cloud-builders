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
	"testing"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/testservices"
)

func TestGetDigest(t *testing.T) {
	ctx := context.Background()
	image := "my-image:1.0.0"
	gs := &testservices.TestGcloud{
		ContainerImagesDescribeResp: "sha256:foobar",
		ContainerImagesDescribeErr:  nil,
	}

	want := "sha256:foobar"

	if got, err := GetDigest(ctx, image, gs); got != want || err != nil {
		t.Errorf("GetDigest(ctx, %s, gs) = %s, %v; want %s, <nil>", image, got, err, want)
	}
}

func TestGetDigestErrors(t *testing.T) {
	ctx := context.Background()
	image := "my-image"
	gs := &testservices.TestGcloud{
		ContainerImagesDescribeResp: "",
		ContainerImagesDescribeErr:  fmt.Errorf("failed to describe container image"),
	}

	if got, err := GetDigest(ctx, image, gs); got != "" || err == nil {
		t.Errorf("GetDigest(ctx, %s, gs) = %s, %v; want \"\", error", image, got, err)
	}
}
