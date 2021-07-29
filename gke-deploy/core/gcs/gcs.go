// Package gcs contains logic related to Google Cloud Storage.
package gcs

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
)

const (
	timeoutErr = "GCS timeout"
)

var (
	defaultTimeout = 1 * time.Hour
	defaultDelay   = 2 * time.Second
	errGCSTimeout  = errors.New(timeoutErr)
)

type GCS struct {
	Timeout    time.Duration
	Retries    int
	Delay      time.Duration
	GcsService services.GcsService
}

// Download copies file(s) from GCS. <dst> should be a directory, not a path to a file.
func (s *GCS) Download(ctx context.Context, src, dst string, recursive bool) error {
	return s.copyWithRetry(ctx, src, dst, recursive)
}

// Upload copies file(s) to GCS.
func (s *GCS) Upload(ctx context.Context, src, dst string) error {
	return s.copyWithRetry(ctx, src, dst, false)
}

// copyWithRetry is responsible for trying (and retrying) to call copyWithTimeout()
// with appropriate retry backoff.
func (s *GCS) copyWithRetry(ctx context.Context, src, dst string, recursive bool) error {
	var err error
	delay := s.Delay
	if delay == 0 {
		delay = defaultDelay
	}
	for retryNum := 0; retryNum <= s.Retries; retryNum++ {
		if retryNum > 0 {
			time.Sleep(delay)
		}
		timeout := s.timeout()
		e := s.copyWithTimeout(ctx, src, dst, recursive, timeout)
		if e != nil {
			err = e
			if strings.Contains(err.Error(), "AccessDeniedException") {
				return err
			}
			continue
		}
		return nil
	}

	return err

}

// copyWithTimeout calls GcsService to move the files. The call will time out if it
// takes too long.
func (s *GCS) copyWithTimeout(ctx context.Context, src, dst string, recursive bool, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	err := s.GcsService.Copy(ctx, src, dst, recursive)
	if ctx.Err() == context.DeadlineExceeded {
		return errGCSTimeout
	}
	return err

}

//timeout returns the GCS timeout that will be used to call GcsService.
func (s *GCS) timeout() time.Duration {
	if s.Timeout == 0 {
		return defaultTimeout
	}
	return s.Timeout
}
