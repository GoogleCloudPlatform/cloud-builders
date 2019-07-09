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
package common

import (
	"reflect"
	"testing"
)

func TestValidateImages(t *testing.T) {
	tests := []struct {
		name string

		images []string
	}{
		{
			name: "No image",

			images: nil,
		},
		{
			name: "One image",

			images: []string{
				"gcr.io/my-project/my-app:1.0.0",
			},
		},
		{
			name: "No duplicates",

			images: []string{
				"gcr.io/my-project/my-app:1.0.0",
				"gcr.io/my-project/my-app2:1.0.0",
				"gcr.io/my-project/my-app3:1.0.0",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateImages(tc.images); err != nil {
				t.Errorf("ValidateImages(%v) = %v; want <nil>", tc.images, err)
			}
		})
	}
}

func TestValidateImagesErrors(t *testing.T) {
	tests := []struct {
		name string

		images []string
	}{
		{
			name: "Duplicate image name",

			images: []string{
				"gcr.io/my-project/my-app:1.0.0",
				"gcr.io/my-project/my-app:1.0.1",
			},
		},
		{
			name: "No tag",

			images: []string{
				"gcr.io/my-project/my-app",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateImages(tc.images); err == nil {
				t.Errorf("ValidateImages(%v) = %v; want err", tc.images, err)
			}
		})
	}
}

func TestCreateLabelsMap(t *testing.T) {
	tests := []struct {
		name string

		labels []string

		want map[string]string
	}{
		{
			name: "Normal case",

			labels: []string{
				"a=b",
				"c=d",
				"d=f",
			},

			want: map[string]string{
				"a": "b",
				"c": "d",
				"d": "f",
			},
		},
		{
			name: "No labels",

			labels: []string{},

			want: map[string]string{},
		},
		{
			name: "Trailing comma",

			labels: []string{
				"a=b",
				"c=d",
				"d=f,",
			},

			want: map[string]string{
				"a": "b",
				"c": "d",
				"d": "f",
			},
		},
		{
			name: "Leading comma",

			labels: []string{
				",a=b",
				"c=d",
				"d=f",
			},

			want: map[string]string{
				"a": "b",
				"c": "d",
				"d": "f",
			},
		},
		{
			name: "Trailing whitespace",

			labels: []string{
				"a=b",
				"c=d",
				"d=f    ",
			},

			want: map[string]string{
				"a": "b",
				"c": "d",
				"d": "f",
			},
		},
		{
			name: "Leading whitespace",

			labels: []string{
				"     a=b",
				"c=d",
				"d=f",
			},

			want: map[string]string{
				"a": "b",
				"c": "d",
				"d": "f",
			},
		},
		{
			name: "Handles whitespace",

			labels: []string{
				" \n a = b  \n\n",
				"\t c = \nd ",
				"d =  f\t",
			},

			want: map[string]string{
				"a": "b",
				"c": "d",
				"d": "f",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := CreateLabelsMap(tc.labels); !reflect.DeepEqual(got, tc.want) || err != nil {
				t.Errorf("CreateLabelsMap(%s) = %v, %v; want %v, <nil>", tc.labels, got, err, tc.want)
			}
		})
	}
}

func TestCreateLabelsMapErrors(t *testing.T) {
	tests := []struct {
		name string

		labels []string
	}{
		{
			name: "No =",

			labels: []string{
				"a",
			},
		},
		{
			name: "More than one =",

			labels: []string{
				"a=b=",
			},
		},
		{
			name: "No key",

			labels: []string{
				"=b",
			},
		},
		{
			name: "No value",

			labels: []string{
				"a=",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := CreateLabelsMap(tc.labels); got != nil || err == nil {
				t.Errorf("CreateLabelsMap(%s) = %v, %v; want <nil>, err", tc.labels, got, err)
			}
		})
	}
}
