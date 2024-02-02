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
			"d=f,",
		},

		want: []applicationsv1beta1.Link{{
			Description: "d",
			URL:         "f",
		}},
	}, {
		name: "Leading comma",

		keyValues: []string{
			",a=b",
		},

		want: []applicationsv1beta1.Link{{
			Description: "a",
			URL:         "b",
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
	}, {
		name: "More than one =",

		keyValues: []string{
			"a=b=",
		},

		want: []applicationsv1beta1.Link{{
			Description: "a",
			URL:         "b=",
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

func TestUseGsutil(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		output   string
		want     bool
	}{{
		name:     "both start with gs://",
		filename: "gs://foo/bar",
		output:   "gs://bar/baz",
		want:     true,
	}, {
		name:     "filename starts with gs://",
		filename: "gs://foo/bar",
		output:   "bar/baz",
		want:     true,
	}, {
		name:     "output start with gs://",
		filename: "foo/bar",
		output:   "gs://bar/baz",
		want:     true,
	}, {
		name:     "netther starts with gs://",
		filename: "foo/bar",
		output:   "bar/baz",
		want:     false,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := UseGsutil(tc.filename, tc.output); got != tc.want {
				t.Errorf("UseGsutil(%s, %s) = got: %v, want: %v", tc.filename, tc.output, got, tc.want)
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
			"d=f,",
		},

		want: map[string]string{
			"d": "f",
		},
	}, {
		name: "Leading comma",

		keyValues: []string{
			",a=b",
		},

		want: map[string]string{
			"a": "b",
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
	}, {
		name: "More than one =",

		keyValues: []string{
			"a=b=",
		},

		want: map[string]string{
			"a": "b=",
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
