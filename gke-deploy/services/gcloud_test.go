package services

import (
	"reflect"
	"testing"

	"google.golang.org/api/container/v1"
)

func TestGetCredentialsArgs(t *testing.T) {
	tests := []struct {
		name          string
		useInternalIP bool
		want          []string
	}{
		{
			name:          "external endpoint",
			useInternalIP: false,
			want: []string{
				"container",
				"clusters",
				"get-credentials",
				"test-cluster",
				"--zone=us-east1-b",
				"--project=my-project",
				"--quiet",
			},
		},
		{
			name:          "internal endpoint",
			useInternalIP: true,
			want: []string{
				"container",
				"clusters",
				"get-credentials",
				"test-cluster",
				"--zone=us-east1-b",
				"--project=my-project",
				"--internal-ip",
				"--quiet",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getCredentialsArgs("test-cluster", "us-east1-b", "my-project", tc.useInternalIP)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("getCredentialsArgs(...) = %v; want %v", got, tc.want)
			}
		})
	}
}

func TestClusterEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		cluster       *container.Cluster
		useInternalIP bool
		want          string
		wantErr       bool
	}{
		{
			name: "external endpoint",
			cluster: &container.Cluster{
				Name:     "test-cluster",
				Endpoint: "34.0.0.1",
				PrivateClusterConfig: &container.PrivateClusterConfig{
					PrivateEndpoint: "10.0.0.2",
				},
			},
			want: "34.0.0.1",
		},
		{
			name:          "internal endpoint",
			useInternalIP: true,
			cluster: &container.Cluster{
				Name:     "test-cluster",
				Endpoint: "34.0.0.1",
				PrivateClusterConfig: &container.PrivateClusterConfig{
					PrivateEndpoint: "10.0.0.2",
				},
			},
			want: "10.0.0.2",
		},
		{
			name:          "missing private endpoint",
			useInternalIP: true,
			cluster: &container.Cluster{
				Name:     "test-cluster",
				Endpoint: "34.0.0.1",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := clusterEndpoint(tc.cluster, tc.useInternalIP)
			if got != tc.want || (err != nil) != tc.wantErr {
				t.Errorf("clusterEndpoint(...) = %q, %v; want %q, err=%v", got, err, tc.want, tc.wantErr)
			}
		})
	}
}
