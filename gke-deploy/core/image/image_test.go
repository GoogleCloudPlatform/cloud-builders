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
package image

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/testservices"
	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/google/go-containerregistry/pkg/name"
)

func TestParseImages(t *testing.T) {
	image := "gcr.io/my-project/my-image:1.0.0"
	image2 := "gcr.io/my-project/my-image-2@sha256:929665b8eb2bb286535d29cd73c71808d7e1ad830046333f6cf0ce497996eb79"
	image3 := "gcr.io/my-project/my-image-3"

	tests := []struct {
		name string

		images []string

		want []name.Reference
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

			want: []name.Reference{
				newImageWithTag(t, image),
				newImageWithDigest(t, image2),
				newImageWithTag(t, image3),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := ParseReferences(tc.images); !reflect.DeepEqual(got, tc.want) || err != nil {
				t.Errorf("ParseReferences(%v) = %v, %v; want %v, <nil>", tc.images, got, err, tc.want)
			}
		})
	}
}

func TestParseImagesErrors(t *testing.T) {
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
			if got, err := ParseReferences(tc.images); err == nil {
				t.Errorf("ParseReferences(%v) = %v, <nil>; want <nil>, err", tc.images, got)
			}
		})
	}
}

func TestName(t *testing.T) {
	tests := []struct {
		name string

		image name.Reference

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
			if got := Name(tc.image); got != tc.want {
				t.Errorf("Name(%v) = %s; want %s", tc.image, got, tc.want)
			}
		})
	}
}

func TestResolveDigest(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		image name.Reference
		rs    *testservices.TestRemote

		want string
	}{
		{
			name: "Get digest from remote image",

			image: newImageWithTag(t, "my-image:1.0.0"),
			rs: &testservices.TestRemote{
				ImageResp: &testservices.TestImage{
					Hash: v1.Hash{
						Algorithm: "sha256",
						Hex:       "foobar",
					},
					Err: nil,
				},
				ImageErr: nil,
			},

			want: "sha256:foobar",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := ResolveDigest(ctx, tc.image, tc.rs); got != tc.want || err != nil {
				t.Errorf("ResolveDigest(ctx, %v, rs) = %s, %v; want %s, <nil>", tc.image, got, err, tc.want)
			}
		})
	}
}

func TestResolveDigestErrors(t *testing.T) {
	ctx := context.Background()
	image := newImageWithTag(t, "my-image:1.0.0")

	tests := []struct {
		name string

		image name.Reference
		rs    *testservices.TestRemote
	}{
		{
			name: "Fail to get remote image",

			image: image,
			rs: &testservices.TestRemote{
				ImageResp: nil,
				ImageErr:  fmt.Errorf("failed to get remote image"),
			},
		},
		{
			name: "Fail to get digest from remote image",

			image: image,
			rs: &testservices.TestRemote{
				ImageResp: &testservices.TestImage{
					Hash: v1.Hash{},
					Err:  fmt.Errorf("failed to get digest"),
				},
				ImageErr: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := ResolveDigest(ctx, image, tc.rs); got != "" || err == nil {
				t.Errorf("ResolveDigest(ctx, %v, rs) = %s, %v; want \"\", error", image, got, err)
			}
		})
	}
}

func newImageWithTag(t *testing.T, image string) name.Reference {
	ref, err := name.NewTag(image)
	if err != nil {
		t.Fatalf("failed to create image with tag: %v", err)
	}
	return ref
}

func newImageWithDigest(t *testing.T, image string) name.Reference {
	ref, err := name.NewDigest(image)
	if err != nil {
		t.Fatalf("failed to create image with digest: %v", err)
	}
	return ref
}
