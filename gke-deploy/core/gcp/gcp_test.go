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
