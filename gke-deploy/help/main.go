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
// Package main runs the doc generation code.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var dir string
var root = &cobra.Command{
	Use:   "gendoc",
	Short: "Generate gke-deploy's help docs",
	Args:  cobra.NoArgs,
	Run: func(*cobra.Command, []string) {
		if err := doc.GenMarkdownTree(cmd.NewCommand(), dir); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	root.Flags().StringVarP(&dir, "dir", "d", ".", "Path to directory in which to generate docs")
}

func main() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
