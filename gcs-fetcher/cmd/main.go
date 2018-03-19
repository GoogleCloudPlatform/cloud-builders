package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/GoogleCloudPlatform/cloud-builders/gcs-fetcher/pkg/fetcher"

	"cloud.google.com/go/storage"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const (
	stagingFolder = ".download/"
	userAgent     = "gae-fetcher"

	defaultWorkers = 200
	defaultRetries = 3
	defaultBackoff = 100 * time.Millisecond
)

var (
	sourceType = flag.String("type", "", "Type of source to fetch; one of Archive or Manifest")
	location   = flag.String("location", "", "Location of source to fetch; in the form gs://bucket/path/to/object#generation")

	destDir     = flag.String("dest_dir", "", "The root where to write the files.")
	workerCount = flag.Int("workers", defaultWorkers, "The number of files to fetch in parallel.")
	verbose     = flag.Bool("verbose", false, "If true, additional output is logged.")
	retries     = flag.Int("retries", defaultRetries, "Number of times to retry a failed GCS download.")
	backoff     = flag.Duration("backoff", defaultBackoff, "Time to wait when retrying, will be doubled on each retry.")
	timeoutGCS  = flag.Bool("timeout_gcs", true, "If true, a timeout will be used to avoid GCS longtails.")
)

func main() {
	flag.Parse()

	if *location == "" || *sourceType == "" {
		glog.Fatal("Must specify --location and --type")
	}

	ctx := context.Background()
	hc, err := buildHTTPClient(ctx)
	if err != nil {
		glog.Info(err)
		os.Exit(2)
	}

	client, err := storage.NewClient(ctx, option.WithHTTPClient(hc), option.WithUserAgent(userAgent))
	if err != nil {
		glog.Infof("Failed to create new GCS client: %v", err)
		os.Exit(2)
	}

	bucket, object, generation, err := fetcher.ParseBucketObject(*location)
	if err != nil {
		glog.Fatalf("Failed to parse --location: %v", err)
	}

	gcs := &fetcher.GCSFetcher{
		GCS:         realGCS{client},
		OS:          realOS{},
		DestDir:     *destDir,
		StagingDir:  filepath.Join(*destDir, stagingFolder),
		CreatedDirs: map[string]bool{},
		Bucket:      bucket,
		Object:      object,
		Generation:  generation,
		TimeoutGCS:  *timeoutGCS,
		WorkerCount: *workerCount,
		Retries:     *retries,
		Backoff:     *backoff,
		SourceType:  *sourceType,
		Verbose:     *verbose,
	}
	if err := gcs.Fetch(ctx); err != nil {
		glog.Fatal(err)
	}
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

func (gp realGCS) NewReader(ctx context.Context, bucket, object string) (io.ReadCloser, error) {
	return gp.client.Bucket(bucket).Object(object).NewReader(ctx)
}

// realOS merely wraps the os package implementations.
type realOS struct{}

func (realOS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (realOS) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (realOS) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (realOS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (realOS) Open(name string) (*os.File, error) {
	return os.Open(name)
}
func (realOS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
