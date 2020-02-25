package crd

import (
	"context"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/testservices"
)

func TestEnsureInstallApplicationCRD(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		ks services.KubectlService
	}{{
		name: "Application CRD already installed",

		ks: &testservices.TestKubectl{
			GetResponse: map[string]map[string][]testservices.GetResponse{
				applicationCRDName: {
					"": {
						{
							Res: "application crd yaml not empty",
							Err: nil,
						},
					},
				},
			},
		},
	}, {
		name: "Install Application CRD",

		ks: &testservices.TestKubectl{
			GetResponse: map[string]map[string][]testservices.GetResponse{
				applicationCRDName: {
					"": {
						{
							Res: "",
							Err: nil,
						},
					},
				},
			},
			ApplyResponse: map[string][]error{
				applicationCRDInstallURI: {nil},
			},
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := EnsureInstallApplicationCRD(ctx, tc.ks); err != nil {
				t.Errorf("EnsureInstallApplicationCRD(ctx, ks) = %v; want <nil>", err)
			}
		})
	}
}

func TestEnsureInstallApplicationCRDErrors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		ks services.KubectlService
	}{{
		name: "Failed to check if Application CRD exists",

		ks: &testservices.TestKubectl{
			GetResponse: map[string]map[string][]testservices.GetResponse{
				applicationCRDName: {
					"": {
						{
							Res: "",
							Err: fmt.Errorf("failed to get Application CRD"),
						},
					},
				},
			},
		},
	}, {
		name: "Failed to install Application CRD",

		ks: &testservices.TestKubectl{
			GetResponse: map[string]map[string][]testservices.GetResponse{
				applicationCRDName: {
					"": {
						{
							Res: "",
							Err: nil,
						},
					},
				},
			},
			ApplyResponse: map[string][]error{
				applicationCRDInstallURI: {fmt.Errorf("failed to install")},
			},
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := EnsureInstallApplicationCRD(ctx, tc.ks); err == nil {
				t.Errorf("EnsureInstallApplicationCRD(ctx, ks) = <nil>; want error")
			}
		})
	}
}
