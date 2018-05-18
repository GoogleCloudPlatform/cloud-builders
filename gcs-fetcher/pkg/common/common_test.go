/*
Copyright 2018 Google, Inc. All rights reserved.

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
package common

import (
	"testing"
)

func TestParseBucketObject(t *testing.T) {
	for _, c := range []struct {
		uri     string
		bucket  string
		object  string
		wantErr bool
	}{{
		uri:    "https://storage.googleapis.com/staging.appid.appspot.com/abc123",
		bucket: "staging.appid.appspot.com",
		object: "abc123",
	}, {
		uri:    "https://storage.googleapis.com/some-bucket.google.com.a.appspot.com/some/path/to/file",
		bucket: "some-bucket.google.com.a.appspot.com",
		object: "some/path/to/file",
	}, {
		uri:    "https://storage.googleapis.com/some-bucket/abc123",
		bucket: "some-bucket",
		object: "abc123",
	}, {
		uri:     "https://storage.googleapis.com/too-short",
		wantErr: true,
	}, {
		uri:     "https://incorrect-domain.com/some-bucket.google.com.a.appspot.com/some/path",
		wantErr: true,
	}, {
		uri:    "gs://my-bucket/manifest-20171004T175409.json",
		bucket: "my-bucket",
		object: "manifest-20171004T175409.json",
	}, {
		uri:    "gs://staging.appid.appspot.com/abc123",
		bucket: "staging.appid.appspot.com",
		object: "abc123",
	}, {
		uri:    "gs://some-bucket.google.com.a.appspot.com/some/path/to/file",
		bucket: "some-bucket.google.com.a.appspot.com",
		object: "some/path/to/file",
	}, {
		uri:    "gs://some-bucket/abc123",
		bucket: "some-bucket",
		object: "abc123",
	}, {
		uri:    "http://storage.googleapis.com/my-bucket/test-memchache/server.js",
		bucket: "my-bucket",
		object: "test-memchache/server.js",
	}, {
		uri:     "gs://too-short",
		wantErr: true,
	}, {
		uri:     "some-bucket/some/path",
		wantErr: true,
	}} {
		bucket, object, _, err := ParseBucketObject(c.uri)
		if (err != nil) != c.wantErr {
			t.Errorf("ParseBucketObject(%q): got %v, wantErr = %t", c.uri, err, c.wantErr)
		}
		if err == nil {
			if bucket != c.bucket || object != c.object {
				t.Errorf("parseBucketObject(%q) = (%q, %q); want (%q, %q)", c.uri, bucket, object, c.bucket, c.object)
			}
		}
	}
}
