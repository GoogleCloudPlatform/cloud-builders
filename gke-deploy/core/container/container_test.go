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
	"reflect"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/testservices"
)

func TestParseImages(t *testing.T) {
	ctx := context.Background()

	image := "gcr.io/my-project/my-image:1.0.0"
	image2 := "gcr.io/my-project/my-image-2@sha256:929665b8eb2bb286535d29cd73c71808d7e1ad830046333f6cf0ce497996eb79"
	image3 := "gcr.io/my-project/my-image-3"

	tests := []struct {
		name string

		images []string

		want []Image
	}{
		{
			name: "No images",

			images: []string{},

			want: nil,
		},
		{
			name: "Success with multiple images",

			images: []string{
				image,
				image2,
				image3,
			},

			want: []Image{
				newImageWithTag(t, image),
				newImageWithDigest(t, image2),
				newImageWithTag(t, image3),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := ParseImages(ctx, tc.images); !reflect.DeepEqual(got, tc.want) || err != nil {
				t.Errorf("ParseImages(ctx, %v) = %v, %v; want %v, <nil>", tc.images, got, err, tc.want)
			}
		})
	}
}

func TestParseImagesErrors(t *testing.T) {
	ctx := context.Background()

	image := "gcr.io/my-project/my-image:1.0.0"
	image2 := "gcr.io/my-project/my-image@sha256:929665b8eb2bb286535d29cd73c71808d7e1ad830046333f6cf0ce497996eb79"
	image3 := "gcr.io/my-project/my-image"

	tests := []struct {
		name string

		images []string
	}{
		{
			name: "Duplicate image name 1",

			images: []string{
				image,
				image2,
			},
		},
		{
			name: "Duplicate image name 2",

			images: []string{
				image,
				image3,
			},
		},
		{
			name: "Duplicate image name 3",

			images: []string{
				image2,
				image3,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := ParseImages(ctx, tc.images); err == nil {
				t.Errorf("ParseImages(ctx, %v) = %v, <nil>; want <nil>, err", tc.images, got)
			}
		})
	}
}

func TestGetName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		image Image

		want string
	}{
		{
			name: "Get name from image with tag",

			image: newImageWithTag(t, "gcr.io/my-project/my-image:1.0.0"),

			want: "gcr.io/my-project/my-image",
		},
		{
			name: "Get name from image with digest",

			image: newImageWithDigest(t, "gcr.io/my-project/my-image-2@sha256:929665b8eb2bb286535d29cd73c71808d7e1ad830046333f6cf0ce497996eb79"),

			want: "gcr.io/my-project/my-image-2",
		},
		{
			name: "Get name from image with latest tag",

			image: newImageWithTag(t, "gcr.io/my-project/my-image-3:latest"),

			want: "gcr.io/my-project/my-image-3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := GetName(ctx, tc.image); got != tc.want || err != nil {
				t.Errorf("GetName(ctx, %v) = %s, %v; want %s, <nil>", tc.image, got, err, tc.want)
			}
		})
	}
}

func TestGetDigest(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		image Image
		gs    *testservices.TestGcloud

		want string
	}{
		{
			name: "Get digest from image with tag",

			image: newImageWithTag(t, "my-image:1.0.0"),
			gs: &testservices.TestGcloud{
				ContainerImagesDescribeResp: "sha256:foobar",
				ContainerImagesDescribeErr:  nil,
			},

			want: "sha256:foobar",
		},
		{
			name:  "Get digest from image with digest",
			image: newImageWithDigest(t, "my-image@sha256:929665b8eb2bb286535d29cd73c71808d7e1ad830046333f6cf0ce497996eb79"),
			gs:    nil,

			want: "sha256:929665b8eb2bb286535d29cd73c71808d7e1ad830046333f6cf0ce497996eb79",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := GetDigest(ctx, tc.image, tc.gs); got != tc.want || err != nil {
				t.Errorf("GetDigest(ctx, %v, gs) = %s, %v; want %s, <nil>", tc.image, got, err, tc.want)
			}
		})
	}
}

func TestGetDigestErrors(t *testing.T) {
	ctx := context.Background()
	image := newImageWithTag(t, "my-image:1.0.0")
	gs := &testservices.TestGcloud{
		ContainerImagesDescribeResp: "",
		ContainerImagesDescribeErr:  fmt.Errorf("failed to describe container image"),
	}

	if got, err := GetDigest(ctx, image, gs); got != "" || err == nil {
		t.Errorf("GetDigest(ctx, %v, gs) = %s, %v; want \"\", error", image, got, err)
	}
}

func newImageWithTag(t *testing.T, image string) Image {
	ref, err := name.NewTag(image)
	if err != nil {
		t.Fatalf("failed to create image with tag: %v", err)
	}
	return ref
}

func newImageWithDigest(t *testing.T, image string) Image {
	ref, err := name.NewDigest(image)
	if err != nil {
		t.Fatalf("failed to create image with digest: %v", err)
	}
	return ref
}
