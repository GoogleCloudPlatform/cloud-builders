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
// Package gcs contains logic related to Google Cloud Storage.
package gcs

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
)

const TimeoutErr = "GCS timeout"

var (
	defaultTimeout = 1 * time.Hour
	errGCSTimeout  = errors.New(TimeoutErr)
)

type GCS struct {
	Timeout    time.Duration
	Retries    int
	Delay      time.Duration
	GcsService services.GcsService
}

// Download copies file(s) from GCS. <dst> should be a directory, not a path to a file
func (s *GCS) Download(ctx context.Context, src, dst string, recursive bool) error {
	log.Printf("Downlaoding file(s) from GCS. Source: %s  Destination: %s.\n", src, dst)
	return s.copyWithRetry(ctx, src, dst, recursive)
}

// Upload copies file(s) to GCS.
func (s *GCS) Upload(ctx context.Context, src, dst string) error {
	log.Printf("Uploading file(s) to GCS.  Source: %s  Destination: %s.\n", src, dst)
	return s.copyWithRetry(ctx, src, dst, false)

}

// copyWithRetry is responsible for trying (and retrying) to call copyWithTimeout()
// with appropriate retry backoff.
func (s *GCS) copyWithRetry(ctx context.Context, src, dst string, recursive bool) error {
	var err error
	delay := s.Delay
	for retryNum := 0; retryNum <= s.Retries; retryNum++ {

		if retryNum > 0 {
			time.Sleep(delay)
			delay *= 2
		}

		started := time.Now()
		timeout := s.timeout()
		e := s.copyWithTimeout(ctx, src, dst, recursive, timeout)
		if e != nil {
			err = e
			log.Printf("Started copying at %v and failed at %v.", started, time.Now())
			continue
		}
		log.Printf("Started copying at %v and finished at %v.", started, time.Now())
		break
	}

	return err

}

// copyWithTimeout calls GcsService to move the files. The call will time out if it
// takes too long.
func (s *GCS) copyWithTimeout(ctx context.Context, src, dst string, recursive bool, timeout time.Duration) error {
	status := make(chan error, 1)
	log.Printf("Operation will time out in %f seconds", timeout.Seconds())
	go func() {
		status <- s.GcsService.Copy(ctx, src, dst, recursive)
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return errGCSTimeout
		}
		return ctx.Err()
	case <-time.After(timeout):
		return errGCSTimeout
	case err := <-status:
		return err
	}

}

//timeout returns the GCS timeout that will be used to call GcsService.
func (s *GCS) timeout() time.Duration {
	if int64(s.Timeout) == 0 {
		return defaultTimeout
	}
	return s.Timeout
}
