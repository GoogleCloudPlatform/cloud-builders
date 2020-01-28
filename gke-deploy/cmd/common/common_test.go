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

	applicationsv1beta1 "github.com/kubernetes-sigs/application/pkg/apis/app/v1beta1"
)

func TestCreateApplicationLinksListFromEqualDelimitedStrings(t *testing.T) {
	tests := []struct {
		name string

		keyValues []string

		want []applicationsv1beta1.Link
	}{{
		name: "Normal case",

		keyValues: []string{
			"My description=https://link.com/",
			"Description2=not really a link",
		},

		want: []applicationsv1beta1.Link{{
			Description: "My description",
			URL:         "https://link.com/",
		}, {
			Description: "Description2",
			URL:         "not really a link",
		}},
	}, {
		name: "No keyValues",

		keyValues: []string{},

		want: nil,
	}, {
		name: "Trailing comma",

		keyValues: []string{
			"a=b",
			"c=d",
			"d=f,",
		},

		want: []applicationsv1beta1.Link{{
			Description: "a",
			URL:         "b",
		}, {
			Description: "c",
			URL:         "d",
		}, {
			Description: "d",
			URL:         "f",
		}},
	}, {
		name: "Leading comma",

		keyValues: []string{
			",a=b",
			"c=d",
			"d=f",
		},

		want: []applicationsv1beta1.Link{{
			Description: "a",
			URL:         "b",
		}, {
			Description: "c",
			URL:         "d",
		}, {
			Description: "d",
			URL:         "f",
		}},
	}, {
		name: "Trailing whitespace",

		keyValues: []string{
			"a=b",
			"c=d",
			"d=f    ",
		},

		want: []applicationsv1beta1.Link{{
			Description: "a",
			URL:         "b",
		}, {
			Description: "c",
			URL:         "d",
		}, {
			Description: "d",
			URL:         "f",
		}},
	}, {
		name: "Leading whitespace",

		keyValues: []string{
			"     a=b",
			"c=d",
			"d=f",
		},

		want: []applicationsv1beta1.Link{{
			Description: "a",
			URL:         "b",
		}, {
			Description: "c",
			URL:         "d",
		}, {
			Description: "d",
			URL:         "f",
		}},
	}, {
		name: "Handles whitespace",

		keyValues: []string{
			" \n a = b  \n\n",
			"\t c = \nd ",
			"d =  f\t",
		},

		want: []applicationsv1beta1.Link{{
			Description: "a",
			URL:         "b",
		}, {
			Description: "c",
			URL:         "d",
		}, {
			Description: "d",
			URL:         "f",
		}},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := CreateApplicationLinksListFromEqualDelimitedStrings(tc.keyValues); !reflect.DeepEqual(got, tc.want) || err != nil {
				t.Errorf("CreateApplicationLinksListFromEqualDelimitedStrings(%s) = %v, %v; want %v, <nil>", tc.keyValues, got, err, tc.want)
			}
		})
	}
}

func TestCreateApplicationLinksListFromEqualDelimitedStringsErrors(t *testing.T) {
	tests := []struct {
		name string

		keyValues []string
	}{{
		name: "No =",

		keyValues: []string{
			"a",
		},
	}, {
		name: "More than one =",

		keyValues: []string{
			"a=b=",
		},
	}, {
		name: "No key",

		keyValues: []string{
			"=b",
		},
	}, {
		name: "No value",

		keyValues: []string{
			"a=",
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := CreateApplicationLinksListFromEqualDelimitedStrings(tc.keyValues); got != nil || err == nil {
				t.Errorf("CreateApplicationLinksListFromEqualDelimitedStrings(%s) = %v, %v; want <nil>, err", tc.keyValues, got, err)
			}
		})
	}
}

func TestCreateMapFromEqualDelimitedStrings(t *testing.T) {
	tests := []struct {
		name string

		keyValues []string

		want map[string]string
	}{{
		name: "Normal case",

		keyValues: []string{
			"a=b",
			"c=d",
			"d=f",
		},

		want: map[string]string{
			"a": "b",
			"c": "d",
			"d": "f",
		},
	}, {
		name: "No keyValues",

		keyValues: []string{},

		want: map[string]string{},
	}, {
		name: "Trailing comma",

		keyValues: []string{
			"a=b",
			"c=d",
			"d=f,",
		},

		want: map[string]string{
			"a": "b",
			"c": "d",
			"d": "f",
		},
	}, {
		name: "Leading comma",

		keyValues: []string{
			",a=b",
			"c=d",
			"d=f",
		},

		want: map[string]string{
			"a": "b",
			"c": "d",
			"d": "f",
		},
	}, {
		name: "Trailing whitespace",

		keyValues: []string{
			"a=b",
			"c=d",
			"d=f    ",
		},

		want: map[string]string{
			"a": "b",
			"c": "d",
			"d": "f",
		},
	}, {
		name: "Leading whitespace",

		keyValues: []string{
			"     a=b",
			"c=d",
			"d=f",
		},

		want: map[string]string{
			"a": "b",
			"c": "d",
			"d": "f",
		},
	}, {
		name: "Handles whitespace",

		keyValues: []string{
			" \n a = b  \n\n",
			"\t c = \nd ",
			"d =  f\t",
		},

		want: map[string]string{
			"a": "b",
			"c": "d",
			"d": "f",
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := CreateMapFromEqualDelimitedStrings(tc.keyValues); !reflect.DeepEqual(got, tc.want) || err != nil {
				t.Errorf("CreateMapFromEqualDelimitedStrings(%s) = %v, %v; want %v, <nil>", tc.keyValues, got, err, tc.want)
			}
		})
	}
}

func TestCreateMapFromEqualDelimitedStringsErrors(t *testing.T) {
	tests := []struct {
		name string

		keyValues []string
	}{{
		name: "No =",

		keyValues: []string{
			"a",
		},
	}, {
		name: "More than one =",

		keyValues: []string{
			"a=b=",
		},
	}, {
		name: "No key",

		keyValues: []string{
			"=b",
		},
	}, {
		name: "No value",

		keyValues: []string{
			"a=",
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := CreateMapFromEqualDelimitedStrings(tc.keyValues); got != nil || err == nil {
				t.Errorf("CreateMapFromEqualDelimitedStrings(%s) = %v, %v; want <nil>, err", tc.keyValues, got, err)
			}
		})
	}
}
