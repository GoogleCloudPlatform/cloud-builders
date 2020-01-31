package gcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/testservices"
)

func TestGetProject(t *testing.T) {
	ctx := context.Background()
	gs := &testservices.TestGcloud{
		ConfigGetValueResp: "my-project",
		ConfigGetValueErr:  nil,
	}

	want := "my-project"

	if got, err := GetProject(ctx, gs); got != want || err != nil {
		t.Errorf("GetProject(ctx, gs) = %s, %v; want %s, <nil>", got, err, want)
	}
}

func TestGetProjectErrors(t *testing.T) {
	ctx := context.Background()
	gs := &testservices.TestGcloud{
		ConfigGetValueResp: "",
		ConfigGetValueErr:  fmt.Errorf("failed to get project"),
	}

	if got, err := GetProject(ctx, gs); got != "" || err == nil {
		t.Errorf("GetProject(ctx, gs) = %s, %v; want \"\", error", got, err)
	}
}

func TestGetAccount(t *testing.T) {
	ctx := context.Background()
	gs := &testservices.TestGcloud{
		ConfigGetValueResp: "my-account",
		ConfigGetValueErr:  nil,
	}

	want := "my-account"

	if got, err := GetAccount(ctx, gs); got != want || err != nil {
		t.Errorf("GetAccount(ctx, gs) = %s, %v; want %s, <nil>", got, err, want)
	}
}

func TestGetAccountErrors(t *testing.T) {
	ctx := context.Background()
	gs := &testservices.TestGcloud{
		ConfigGetValueResp: "",
		ConfigGetValueErr:  fmt.Errorf("failed to get account"),
	}

	if got, err := GetAccount(ctx, gs); got != "" || err == nil {
		t.Errorf("GetAccount(ctx, gs) = %s, %v; want \"\", error", got, err)
	}
}
