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
package resource

import (
	"fmt"
	"os"

	"github.com/gobuffalo/packr/v2"
)

var (
	box = packr.New("configTemplates", "./templates")

	namespaceTemplateBytes = readConfigTemplate("namespace.yaml")
)

func readConfigTemplate(filename string) []byte {
	in, err := box.Find(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read resource config template %q: %v", filename, err)
		os.Exit(1)
	}
	return in
}
