/*
Copyright 2018 Google, Inc. All rights reserved.

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
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/GoogleCloudPlatform/cloud-builders/gcs-fetcher/pkg/uploader"
)

const userAgent = "gcs-uploader"

var (
	dir          = flag.String("dir", ".", "Directory of files to upload")
	bucket       = flag.String("bucket", "", "GCS bucket to upload files and manifest to")
	manifestFile = flag.String("manifest_file", "", "If specified, manifest file name; otherwise, one will be generated")
	workerCount  = flag.Int("workers", 200, "The number of files to upload in parallel.")
	help         = flag.Bool("help", false, "If true, prints help text and exits.")
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println("Incrementally uploads source files to Google Cloud Storage")
		flag.PrintDefaults()
		return
	}

	if *bucket == "" {
		log.Fatalf("--bucket must be specified")
	}

	ctx := context.Background()
	hc, err := buildHTTPClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client, err := storage.NewClient(ctx, option.WithHTTPClient(hc), option.WithUserAgent(userAgent))
	if err != nil {
		log.Fatalf("Failed to create new GCS client: %v", err)
	}

	u := uploader.GCSUploader{
		GCS:          realGCS{client},
		OS:           realOS{},
		Root:         *dir,
		Bucket:       *bucket,
		ManifestFile: *manifestFile,
		WorkerCount:  *workerCount,
	}

	manifestURL, err := u.Upload(ctx)
	if err != nil {
		log.Fatalf("Failed to upload: %v", err)
	}

	log.Printf("Uploaded manifest: %s", manifestURL)
}

func buildHTTPClient(ctx context.Context) (*http.Client, error) {
	hc, err := google.DefaultClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create default client: %v", err)
	}

	ts, err := google.DefaultTokenSource(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create default token source: %v", err)
	}

	hc.Transport = &oauth2.Transport{
		Base:   http.DefaultTransport,
		Source: ts,
	}

	return hc, nil
}

// realGCS is a wrapper over the GCS client functions.
type realGCS struct {
	client *storage.Client
}

func (gp realGCS) NewWriter(ctx context.Context, bucket, object string) io.WriteCloser {
	return gp.client.Bucket(bucket).Object(object).
		If(storage.Conditions{DoesNotExist: true}). // Skip upload if already exists.
		NewWriter(ctx)
}

// realOS merely wraps the os package implementations.
type realOS struct{}

func (realOS) Walk(root string, fn filepath.WalkFunc) error { return filepath.Walk(root, fn) }
func (realOS) EvalSymlinks(path string) (string, error)     { return filepath.EvalSymlinks(path) }
func (realOS) Stat(path string) (os.FileInfo, error)        { return os.Stat(path) }
