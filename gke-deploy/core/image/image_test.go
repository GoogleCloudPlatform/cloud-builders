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
	"testing"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/testservices"
	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/google/go-containerregistry/pkg/name"
)

func TestName(t *testing.T) {
	tests := []struct {
		name string

		image name.Reference

		want string
	}{{
		name: "Get name from image with tag",

		image: newImageWithTag(t, "gcr.io/my-project/my-image:1.0.0"),

		want: "gcr.io/my-project/my-image",
	}, {
		name: "Get name from image with digest",

		image: newImageWithDigest(t, "gcr.io/my-project/my-image-2@sha256:929665b8eb2bb286535d29cd73c71808d7e1ad830046333f6cf0ce497996eb79"),

		want: "gcr.io/my-project/my-image-2",
	}, {
		name: "Get name from image with latest tag",

		image: newImageWithTag(t, "gcr.io/my-project/my-image-3:latest"),

		want: "gcr.io/my-project/my-image-3",
	}}

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

	image := newImageWithTag(t, "my-image:1.0.0")
	rs := &testservices.TestRemote{
		ImageResp: &testservices.TestImage{
			Hash: v1.Hash{
				Algorithm: "sha256",
				Hex:       "foobar",
			},
			Err: nil,
		},
		ImageErr: nil,
	}

	want := "sha256:foobar"

	if got, err := ResolveDigest(ctx, image, rs); got != want || err != nil {
		t.Errorf("ResolveDigest(ctx, %v, rs) = %s, %v; want %s, <nil>", image, got, err, want)
	}
}

func TestResolveDigestErrors(t *testing.T) {
	ctx := context.Background()
	image := newImageWithTag(t, "my-image:1.0.0")

	tests := []struct {
		name string

		image name.Reference
		rs    *testservices.TestRemote
	}{{
		name: "Fail to get remote image",

		image: image,
		rs: &testservices.TestRemote{
			ImageResp: nil,
			ImageErr:  fmt.Errorf("failed to get remote image"),
		},
	}, {
		name: "Fail to get digest from remote image",

		image: image,
		rs: &testservices.TestRemote{
			ImageResp: &testservices.TestImage{
				Hash: v1.Hash{},
				Err:  fmt.Errorf("failed to get digest"),
			},
			ImageErr: nil,
		},
	}}

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
