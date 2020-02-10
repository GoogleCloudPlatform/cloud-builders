package gcs

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/testservices"
)

const (
	retries    = 2
	dstDir     = ".workspace/"
	timeoutGCS = 1 * time.Second
	delay      = 1 * time.Second

	singleFile = "gs://k8s.yml"
	directory  = "gs://configs"
	nestedDir  = "gs://configs/nested"
	retryFile  = "gs://k8s_retry.yml"

	slowReadFile     = "gs://project/bucket/slow_file.yaml"
	errorFile        = "gs://project/bucket/not_existed"
	accessDeniedFile = "gs://project/not_authorized_bucket"
	notFoundFile     = "gs://project/bucket/not_found"

	expandedK8sConfig    = "workspace/expanded/aggregated-resources.yaml"
	slowUploadConfig     = "workspace/expanded/aggregated-resources-slow.yaml"
	bucketNotFoundConfig = "workspace/expanded/aggregated-resources-not-found.yaml"
	accessDeniedConfig   = "workspace/expanded/aggregated-resources-access-denied.yaml"

	expandedGcsDst     = "gs://project/path/expanded/aggregated-resources.yaml"
	bucketNotFoundDst  = "gs://project/not_found_bucket/expanded/aggregated-resources.yaml"
	accessDeniedBucket = "gs://project/access_denied_bucket/expanded/aggregated-resources.yaml"

	errMsg            = "copy file failed"
	errDenied         = "error is AccessDeniedException: 403"
	errFileNotFound   = "error is CommandException: No URLs matched"
	errBucketNotFound = "error is NotFoundException: 404 The destination bucket does not exist"
)

var counter int

func buildTestGCS(t *testing.T) *GCS {
	t.Helper()

	s := &testservices.TestGcsService{CopyResponse: map[string]func(src, dst string, recursive bool) error{
		singleFile: func(src, dst string, recursive bool) error { return nil },
		directory:  func(src, dst string, recursive bool) error { return nil },
		nestedDir:  func(src, dst string, recursive bool) error { return nil },
		slowReadFile: func(src, dst string, recursive bool) error {
			time.Sleep(3 * time.Second)
			return nil
		},
		retryFile: func(src, dst string, recursive bool) error {
			if counter == 0 {
				counter++
				return errors.New(errMsg)
			}
			return nil
		},

		errorFile:         func(src, dst string, recursive bool) error { return errors.New(errMsg) },
		accessDeniedFile:  func(src, dst string, recursive bool) error { return errors.New(errDenied) },
		notFoundFile:      func(src, dst string, recursive bool) error { return errors.New(errFileNotFound) },
		expandedK8sConfig: func(src, dst string, recursive bool) error { return nil },
		slowUploadConfig: func(src, dst string, recursive bool) error {
			time.Sleep(3 * time.Second)
			return nil
		},
		accessDeniedConfig:   func(src, dst string, recursive bool) error { return errors.New(errDenied) },
		bucketNotFoundConfig: func(src, dst string, recursive bool) error { return errors.New(errBucketNotFound) },
	},
	}

	return &GCS{
		Timeout:    timeoutGCS,
		Retries:    retries,
		Delay:      delay,
		GcsService: s,
	}

}

func TestFetch(t *testing.T) {

	tests := []struct {
		name      string
		src       string
		dst       string
		recursive bool
	}{
		{"download single file", singleFile, dstDir, false},
		{"download files from a directory", directory, dstDir, false},
		{"download files from a directory with nested folders", nestedDir, dstDir, true},
		{"error first then retry", retryFile, dstDir, false},
	}

	gcs := buildTestGCS(t)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := gcs.Download(context.Background(), tc.src, tc.dst, tc.recursive); err != nil {
				t.Errorf("Download(ctx, %s, %s, %t) = %v; want <nil>", tc.src, tc.dst, tc.recursive, err)
			}
		})
	}

}

func TestFetchErrors(t *testing.T) {

	tests := []struct {
		name      string
		src       string
		dst       string
		recursive bool
	}{
		{"download files failed on timeout", slowReadFile, dstDir, false},
		{"download files with failure", errorFile, dstDir, false},
		{"download files with access denied failure", accessDeniedFile, dstDir, false},
		{"download files with not found failure", notFoundFile, dstDir, false},
	}

	gcs := buildTestGCS(t)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := gcs.Download(context.Background(), tc.src, tc.dst, tc.recursive)
			if err == nil {
				t.Errorf("Download(ctx, %s, %s, %t) = <nil>; want error", tc.src, tc.dst, tc.recursive)
				return
			}
			if tc.src == slowReadFile && err.Error() != TimeoutErr {
				t.Errorf("error is %t, want %q error", err, "Timout")
			}
			if tc.src == accessDeniedFile && !strings.Contains(err.Error(), "AccessDeniedException") {
				t.Errorf("error is %t, want %q", err, "AccessDeniedException")
			}
			if tc.src == notFoundFile && !strings.Contains(err.Error(), "CommandException") {
				t.Errorf("error is %t, want %q", err, "CommandException")
			}
		})
	}

}

func TestUpload(t *testing.T) {

	tests := []struct {
		name string
		src  string
		dst  string
	}{
		{"upload expended file", expandedK8sConfig, expandedGcsDst},
	}

	gcs := buildTestGCS(t)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := gcs.Upload(context.Background(), tc.src, tc.dst); err != nil {
				t.Errorf("Upload(ctx, %s, %s) = %v; want <nil>", tc.src, tc.dst, err)
			}
		})
	}

}

func TestUploadErrors(t *testing.T) {

	tests := []struct {
		name string
		src  string
		dst  string
	}{
		{"upload files failed on timeout", slowUploadConfig, expandedGcsDst},
		{"upload files with access denied failure", accessDeniedConfig, accessDeniedBucket},
		{"upload files with bucket not found failure", bucketNotFoundConfig, bucketNotFoundDst},
	}

	gcs := buildTestGCS(t)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := gcs.Upload(context.Background(), tc.src, tc.dst)
			if err == nil {
				t.Errorf("Upload(ctx, %s, %s) = <nil>; want error", tc.src, tc.dst)
				return
			}
			if tc.src == slowUploadConfig && err.Error() != TimeoutErr {
				t.Errorf("error is %t, want %q error", err, "Timout")
			}
			if tc.src == accessDeniedConfig && !strings.Contains(err.Error(), "AccessDeniedException") {
				t.Errorf("error is %t, want %q", err, "AccessDeniedException")
			}
			if tc.src == bucketNotFoundConfig && !strings.Contains(err.Error(), "NotFoundException") {
				t.Errorf("error is %t, want %q", err, "NotFoundException")
			}
		})
	}

}
