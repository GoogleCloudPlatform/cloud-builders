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
package testservices

import (
	"context"
	"fmt"
)

// TestKubectl implements the KubectlService interface.
type TestKubectl struct {
	ApplyResponse map[string]error
	GetResponse   map[string]map[string]*GetResponse
}

// StatResponse represents a response tuple for a Stat function call.
type GetResponse struct {
	Res   []string
	Err   []error
	count int
}

// Apply calls `kubectl apply -f <configs> -n <namespace>`.
func (k *TestKubectl) Apply(ctx context.Context, configs, namespace string) error {
	err, ok := k.ApplyResponse[configs]
	if !ok {
		panic(fmt.Sprintf("ApplyResponse has no response for configs %q", configs))
	}
	return err
}

// Get calls `kubectl get <kind> <name> -n <namespace> --output=<format>`.
func (k *TestKubectl) Get(ctx context.Context, kind, name, namespace, format string) (string, error) {
	resp, ok := k.GetResponse[kind][name]
	if !ok {
		panic(fmt.Sprintf("GetResponse has no response for kind %q and name %q", kind, name))
	}

	defer func() {
		resp.count++
	}()
	return resp.Res[resp.count], resp.Err[resp.count]
}
