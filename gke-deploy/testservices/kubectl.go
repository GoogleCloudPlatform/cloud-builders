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
	ApplyFromStringResponse map[string][]error
	GetResponse             map[string]map[string][]GetResponse
}

// StatResponse represents a response tuple for a Stat function call.
type GetResponse struct {
	Res string
	Err error
}

// ApplyFromString calls `kubectl apply -f - -n <namespace> < ${configString}`.
func (k *TestKubectl) ApplyFromString(configString, namespace string) error {
	errors, ok := k.ApplyFromStringResponse[configString]
	if !ok {
		panic(fmt.Sprintf("ApplyFromStringResponse has no response for configs %q", configString))
	}
	if len(errors) == 0 {
		panic(fmt.Sprintf("ApplyFromStringResponse ran out of responses for configs %q", configString))
	}
	err := errors[0]
	if len(errors) == 1 {
		delete(k.ApplyFromStringResponse, configString)
	} else {
		k.ApplyFromStringResponse[configString] = k.ApplyFromStringResponse[configString][1:]
	}
	return err
}

// Get calls `kubectl get <kind> <name> -n <namespace> --output=<format>`.
func (k *TestKubectl) Get(ctx context.Context, kind, name, namespace, format string, ignoreNotFound bool) (string, error) {
	resp, ok := k.GetResponse[kind][name]
	if !ok {
		panic(fmt.Sprintf("GetResponse has no response for kind %q and name %q", kind, name))
	}

	if len(resp) == 0 {
		panic(fmt.Sprintf("GetResponse ran out of responses for kind %q and name %q after", kind, name))
	}
	res := resp[0].Res
	err := resp[0].Err

	if len(resp) == 1 {
		delete(k.GetResponse[kind], name)
		if len(k.GetResponse[kind]) == 0 {
			delete(k.GetResponse, kind)
		}
	} else {
		k.GetResponse[kind][name] = k.GetResponse[kind][name][1:]
	}
	return res, err
}
