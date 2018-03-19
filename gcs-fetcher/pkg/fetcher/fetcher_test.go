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
package fetcher

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

const (
	maxretries = 3

	successBucket     = "success-bucket"
	sfile1            = "sfile1.js"
	sfile2            = "sfile2.jpg"
	sfile3            = "sfile3"
	goodManifest      = "good-manifest.json"
	malformedManifest = "malformed-manifest.json"

	errorBucket   = "error-bucket"
	efile1        = "efile1"
	efile2        = "efile2"
	efile3        = "efile3"
	errorManifest = "error-manifest.json"
	errorZipfile  = "error-source.zip"

	generation int64 = 12345
)

var (
	zeroTime  = time.Time{}
	errNonNil = errors.New("some-error")

	sfile1Contents = []byte("sfile1-contents-a")
	sfile2Contents = []byte("sfile2-contents-aa")
	sfile3Contents = []byte("sfile3-contents-aaa")

	goodManifestContents = []byte(`{
		"sfile1.js":  {"SourceURL": "gs://success-bucket/sfile1.js", "Sha1Sum": ""},
		"sfile2.jpg": {"SourceURL": "gs://success-bucket/sfile2.jpg", "Sha1Sum": ""},
		"sfile3":     {"SourceURL": "gs://success-bucket/sfile3", "Sha1Sum": ""}
	}`)
	malformedManifestContents = []byte(`{
		"sfile1.js": {"SourceURL": "gs://success-bucket/sfile1.js", "Sha1Sum": ""},
		"sfile2.jpg": {"SourceURL": "gs://succ`)

	errGCSNewReader = fmt.Errorf("instrumented GCS NewReader error")
	errGCSRead      = fmt.Errorf("instrumented GCS Read err")
	errGCSSlowRead  = fmt.Errorf("instrumented GCS slow Read err")
	errRename       = fmt.Errorf("instrumented os.Rename error")
	errChmod        = fmt.Errorf("instrumented os.Chmod error")
	errCreate       = fmt.Errorf("instrumented os.Create error")
	errMkdirAll     = fmt.Errorf("instrumented os.MkdirAll error")
	errOpen         = fmt.Errorf("instrumented os.Open error")
)

type fakeGCSErrorReader struct {
	err   error
	sleep time.Duration
}

func (f fakeGCSErrorReader) Read([]byte) (int, error) {
	time.Sleep(f.sleep)
	return 0, f.err
}

type fakeGCSResponse struct {
	content []byte
	err     error
}

// fakeGCS allows us to simulate errors when interacting with GCS.
type fakeGCS struct {
	t       *testing.T
	objects map[string]fakeGCSResponse
}

func (f *fakeGCS) NewReader(context context.Context, bucket, object string) (io.ReadCloser, error) {
	f.t.Helper()
	name := formatGCSName(bucket, object, generation)

	response, ok := f.objects[name]
	if !ok {
		f.t.Fatalf("no %q in instrumented responses", name)
		return nil, nil
	}

	if response.err == errGCSNewReader {
		return ioutil.NopCloser(bytes.NewReader([]byte(""))), response.err
	}

	if response.err == errGCSRead {
		return ioutil.NopCloser(fakeGCSErrorReader{err: response.err}), nil
	}

	if response.err == errGCSSlowRead {
		return ioutil.NopCloser(fakeGCSErrorReader{sleep: 1 * time.Second}), nil
	}

	if response.err != nil {
		f.t.Fatalf("unexpected error type %v", response.err)
	}

	return ioutil.NopCloser(bytes.NewReader(response.content)), nil
}

// fakeOS raises errors if configures, otherwise simply passes
// through to the normal os package.
type fakeOS struct {
	errorsRename   int
	errorsChmod    int
	errorsCreate   int
	errorsMkdirAll int
	errorsOpen     int
}

func (f *fakeOS) Rename(oldpath, newpath string) error {
	if f.errorsRename > 0 {
		f.errorsRename--
		return errRename
	}
	return os.Rename(oldpath, newpath)
}

func (f *fakeOS) Chmod(name string, mode os.FileMode) error {
	if f.errorsChmod > 0 {
		f.errorsChmod--
		return errChmod
	}
	return os.Chmod(name, mode)
}

func (f *fakeOS) Create(name string) (*os.File, error) {
	if f.errorsCreate > 0 {
		f.errorsCreate--
		return nil, errCreate
	}

	return os.Create(name)
}

func (f *fakeOS) MkdirAll(path string, perm os.FileMode) error {
	if f.errorsMkdirAll > 0 {
		f.errorsMkdirAll--
		return errMkdirAll
	}
	return os.MkdirAll(path, perm)
}

func (f *fakeOS) Open(name string) (*os.File, error) {
	if f.errorsOpen > 0 {
		f.errorsOpen--
		return nil, errOpen
	}
	return os.Open(name)
}

func (*fakeOS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

type testContext struct {
	gf      *GCSFetcher
	gcs     *fakeGCS
	os      *fakeOS
	workDir string
}

func buildManifestTestContext(t *testing.T) (tc *testContext, teardown func()) {
	t.Helper()

	// Set up a temp directory for each test so it's easy to clean up.
	workDir, err := ioutil.TempDir("", "fetcher")
	if err != nil {
		t.Fatal(err)
	}

	os := &fakeOS{}

	gcs := &fakeGCS{
		t: t,
		objects: map[string]fakeGCSResponse{
			formatGCSName(successBucket, sfile1, generation):            {content: sfile1Contents},
			formatGCSName(successBucket, sfile2, generation):            {content: sfile2Contents},
			formatGCSName(successBucket, sfile3, generation):            {content: sfile3Contents},
			formatGCSName(errorBucket, efile1, generation):              {err: errGCSNewReader},
			formatGCSName(errorBucket, efile2, generation):              {err: errGCSRead},
			formatGCSName(errorBucket, efile3, generation):              {err: errGCSSlowRead},
			formatGCSName(successBucket, goodManifest, generation):      {content: goodManifestContents},
			formatGCSName(successBucket, malformedManifest, generation): {content: malformedManifestContents},
			formatGCSName(errorBucket, errorManifest, generation):       {err: errGCSRead},
		},
	}

	gf := &GCSFetcher{
		GCS:         gcs,
		OS:          os,
		DestDir:     workDir,
		StagingDir:  filepath.Join(workDir, ".staging/"),
		CreatedDirs: make(map[string]bool),
		Bucket:      successBucket,
		Object:      goodManifest,
		TimeoutGCS:  true,
		WorkerCount: 2,
		Retries:     maxretries,
	}

	return &testContext{
			workDir: workDir,
			os:      os,
			gcs:     gcs,
			gf:      gf,
		},
		func() {
			if err := os.RemoveAll(workDir); err != nil {
				t.Logf("Failed to remove working dir %q, continuing.", workDir)
			}
		}
}

func TestFetchObjectOnceStoresFile(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()

	j := job{bucket: successBucket, object: sfile1}
	dest := filepath.Join(tc.workDir, "sfile1.tmp")

	result := tc.gf.fetchObjectOnce(context.Background(), j, dest, make(chan struct{}, 1))

	if result.err != nil {
		t.Errorf("fetchObjectOnce() result.err got %v, want nil", result.err)
	}
	if int(result.size) != len(sfile1Contents) {
		t.Errorf("fetchObjectOnce() result.size got %d, want %d", result.size, len(sfile1Contents))
	}

	got, err := ioutil.ReadFile(dest)
	if err != nil {
		t.Fatalf("ReadFile(%v) got %v, want nil", dest, err)
	}
	if !bytes.Equal(got, sfile1Contents) {
		t.Fatalf("ReadFile(%v) got %v, want %v", dest, got, sfile1Contents)
	}
}

func TestFetchObjectOnceFailureModes(t *testing.T) {

	// GCS NewReader failure
	tc, teardown := buildManifestTestContext(t)
	j := job{bucket: errorBucket, object: efile1}
	result := tc.gf.fetchObjectOnce(context.Background(), j, filepath.Join(tc.workDir, "efile1.tmp"), make(chan struct{}, 1))
	if result.err == nil || !strings.HasSuffix(result.err.Error(), errGCSNewReader.Error()) {
		t.Errorf("fetchObjectOnce did not fail correctly, got err=%v, want err=%v", result.err, errGCSNewReader)
	}
	teardown()

	// Failure due to cancellation
	tc, teardown = buildManifestTestContext(t)
	breaker := make(chan struct{}, 1)
	breaker <- struct{}{}
	j = job{bucket: successBucket, object: sfile1}
	result = tc.gf.fetchObjectOnce(context.Background(), j, filepath.Join(tc.workDir, "sfile1.tmp"), breaker)
	if result.err == nil || result.err != errGCSTimeout {
		t.Errorf("fetchObjectOnce did not fail correctly, got err=%v, want err=%v", result.err, errGCSTimeout)
	}
	teardown()

	// os.Create failure
	tc, teardown = buildManifestTestContext(t)
	tc.os.errorsCreate = 1
	j = job{bucket: successBucket, object: sfile1}
	result = tc.gf.fetchObjectOnce(context.Background(), j, filepath.Join(tc.workDir, "sfile1.tmp"), make(chan struct{}, 1))
	if result.err == nil || !strings.HasSuffix(result.err.Error(), errCreate.Error()) {
		t.Errorf("fetchObjectOnce did not fail correctly, got err=%v, want err=%v", result.err, errCreate)
	}
	teardown()

	// GCS Copy failure
	tc, teardown = buildManifestTestContext(t)
	j = job{bucket: errorBucket, object: efile2}
	result = tc.gf.fetchObjectOnce(context.Background(), j, filepath.Join(tc.workDir, "efile2.tmp"), make(chan struct{}, 1))
	if result.err == nil || !strings.HasSuffix(result.err.Error(), errGCSRead.Error()) {
		t.Errorf("fetchObjectOnce did not fail correctly, got err=%v, want err=%v", result.err, errGCSRead)
	}
	teardown()

	// SHA checksum failure
	// TODO(jasonco): Add a SHA checksum failure test
}

func TestFetchObjectOnceWithTimeoutSucceeds(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()

	j := job{bucket: successBucket, object: sfile1}
	timeout := 10 * time.Second
	dest := filepath.Join(tc.workDir, "sfile1.tmp")

	n, err := tc.gf.fetchObjectOnceWithTimeout(context.Background(), j, timeout, dest)
	if err != nil || int(n) != len(sfile1Contents) {
		t.Errorf("fetchObjectOnceWithTimeout() got (%v, %v), want (%v, %v)", n, err, nil, len(sfile1Contents))
	}
}

func TestFetchObjectOnceWithTimeoutFailsOnTimeout(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()

	j := job{bucket: errorBucket, object: efile3} // efile3 is a slow GCS read
	timeout := 100 * time.Millisecond
	dest := filepath.Join(tc.workDir, "efile3.tmp")

	if _, err := tc.gf.fetchObjectOnceWithTimeout(context.Background(), j, timeout, dest); err == nil {
		t.Errorf("fetchObjectOnceWithTimeout() got err=nil, want err=%v", errGCSTimeout)
	}
}

func TestFetchObjectSucceeds(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()

	j := job{bucket: successBucket, object: sfile1, filename: "localfile.txt"}
	report := tc.gf.fetchObject(context.Background(), j)

	if report.job != j {
		t.Errorf("report.job got %v, want %v", report.job, j)
	}
	if !report.success {
		t.Errorf("report.success got false, want true")
	}
	if report.err != nil {
		t.Errorf("report.err got %v, want nil", report.err)
	}
	if report.started == zeroTime {
		t.Errorf("report.started got %v, want report.started > %v", report.started, zeroTime)
	}
	if report.completed == zeroTime {
		t.Errorf("report.completed got %v, want report.completed > %v", report.completed, zeroTime)
	}
	if int(report.size) != len(sfile1Contents) {
		t.Errorf("report.size got %v, want %v", report.size, len(sfile1Contents))
	}
	if report.finalname == "" {
		t.Errorf("report.finalname got empty string, want non-empty string")
	}
	if len(report.attempts) != 1 {
		t.Fatalf("len(report.attempts) got %d, want 1", len(report.attempts))
	}

	attempt := report.attempts[0]
	if attempt.started == zeroTime {
		t.Errorf("attempt.started got %v, want attempt.started > %v", attempt.started, zeroTime)
	}
	if attempt.duration == 0 {
		t.Errorf("attempt.duration got %v, want attempt.duration>0", attempt.duration)
	}
	if attempt.err != nil {
		t.Errorf("attempt.err got %v, want nil", attempt.err)
	}

	got, err := ioutil.ReadFile(report.finalname)
	if err != nil {
		t.Fatalf("ReadFile(%v) got %v, want nil", report.finalname, err)
	}
	if !bytes.Equal(got, sfile1Contents) {
		t.Fatalf("ReadFile(%v) got %v, want %v", report.finalname, got, sfile1Contents)
	}
}

func TestFetchObjectRetriesUntilSuccess(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()
	tc.os.errorsCreate = 1 // first create fails, second succeeds

	j := job{bucket: successBucket, object: sfile1, filename: "localhost.txt"}
	report := tc.gf.fetchObject(context.Background(), j)

	if !report.success {
		t.Errorf("report.success got false, want true")
	}
	if report.err != nil {
		t.Errorf("report.err got %v, want nil", report.err)
	}

	if len(report.attempts) != 2 {
		t.Fatalf("len(report.attempts) got %d, want 2", len(report.attempts))
	}

	attempt1 := report.attempts[0]
	if attempt1.err == nil {
		t.Errorf("attempt.err got %v, want non-nil", attempt1.err)
	}

	attempt2 := report.attempts[1]
	if attempt2.err != nil {
		t.Errorf("attempt.err got %v, want nil", attempt2.err)
	}

	got, err := ioutil.ReadFile(report.finalname)
	if err != nil {
		t.Fatalf("ReadFile(%v) got %v, want nil", report.finalname, err)
	}
	if !bytes.Equal(got, sfile1Contents) {
		t.Fatalf("ReadFile(%v) got %v, want %v", report.finalname, got, sfile1Contents)
	}
}

func TestFetchObjectRetriesMaxTimes(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()
	tc.os.errorsCreate = maxretries + 1 // create continually fails until max reached

	filename := "localfile.txt"
	j := job{bucket: successBucket, object: sfile1, filename: filename}

	report := tc.gf.fetchObject(context.Background(), j)

	if report.success {
		t.Errorf("report.success got true, want false")
	}
	if report.err == nil {
		t.Errorf("report.err got %v, want non-nil", report.err)
	}
	if report.finalname != "" {
		t.Errorf("report.finalname got %v want empty string", report.finalname)
	}
	if len(report.attempts) != maxretries+1 {
		t.Fatalf("len(report.attempts) got %d, want %d", len(report.attempts), maxretries+1)
	}

	last := report.attempts[len(report.attempts)-1]
	if last.err == nil {
		t.Errorf("attempt.err got %v, want non-nil", last.err)
	}

	localfile := filepath.Join(tc.gf.DestDir, filename)
	if _, err := os.Stat(localfile); !os.IsNotExist(err) {
		t.Errorf("file %q exists, want not exists", localfile)
	}
}

func TestFetchObjectRetriesOnFolderCreationError(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()
	tc.os.errorsMkdirAll = 1

	j := job{bucket: successBucket, object: sfile1, filename: "localfile.txt"}
	report := tc.gf.fetchObject(context.Background(), j)

	if !report.success {
		t.Errorf("report.success got false, want true")
	}
	if report.err != nil {
		t.Errorf("report.err got %v, want nil", report.err)
	}

	if len(report.attempts) != 2 {
		t.Fatalf("len(report.attempts) got %d, want 2", len(report.attempts))
	}
	first := report.attempts[0]
	if first.err == nil || !strings.Contains(first.err.Error(), errMkdirAll.Error()) {
		t.Errorf("attempt.err got %v, want Contains(%v)", first.err, errMkdirAll)
	}

	got, err := ioutil.ReadFile(report.finalname)
	if err != nil {
		t.Fatalf("ReadFile(%v) got %v, want nil", report.finalname, err)
	}
	if !bytes.Equal(got, sfile1Contents) {
		t.Fatalf("ReadFile(%v) got %v, want %v", report.finalname, got, sfile1Contents)
	}
}

func TestFetchObjectRetriesOnFetchFail(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()
	tc.os.errorsCreate = 1 // Invoked when fetching the file.

	j := job{bucket: successBucket, object: sfile1, filename: "localfile.txt"}
	report := tc.gf.fetchObject(context.Background(), j)

	if !report.success {
		t.Errorf("report.success got false, want true")
	}
	if report.err != nil {
		t.Errorf("report.err got %v, want nil", report.err)
	}

	if len(report.attempts) != 2 {
		t.Fatalf("len(report.attempts) got %d, want 2", len(report.attempts))
	}
	first := report.attempts[0]
	if first.err == nil || !strings.Contains(first.err.Error(), errCreate.Error()) {
		t.Errorf("attempt.err got %v, want Contains(%v)", first.err, errCreate)
	}

	got, err := ioutil.ReadFile(report.finalname)
	if err != nil {
		t.Fatalf("ReadFile(%v) got %v, want nil", report.finalname, err)
	}
	if !bytes.Equal(got, sfile1Contents) {
		t.Fatalf("ReadFile(%v) got %v, want %v", report.finalname, got, sfile1Contents)
	}
}

func TestFetchObjectRetriesOnRenameFailure(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()
	tc.os.errorsRename = 1

	j := job{bucket: successBucket, object: sfile1, filename: "localfile.txt"}
	report := tc.gf.fetchObject(context.Background(), j)

	if !report.success {
		t.Errorf("report.success got false, want true")
	}
	if report.err != nil {
		t.Errorf("report.err got %v, want nil", report.err)
	}

	if len(report.attempts) != 2 {
		t.Fatalf("len(report.attempts) got %d, want 2", len(report.attempts))
	}
	first := report.attempts[0]
	if first.err == nil || !strings.Contains(first.err.Error(), errRename.Error()) {
		t.Errorf("attempt.err got %v, want Contains(%v)", first.err, errRename)
	}

	got, err := ioutil.ReadFile(report.finalname)
	if err != nil {
		t.Fatalf("ReadFile(%v) got %v, want nil", report.finalname, err)
	}
	if !bytes.Equal(got, sfile1Contents) {
		t.Fatalf("ReadFile(%v) got %v, want %v", report.finalname, got, sfile1Contents)
	}
}

func TestFetchObjectRetriesOnChmodFailure(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()
	tc.os.errorsChmod = 1

	j := job{bucket: successBucket, object: sfile1, filename: "localfile.txt"}
	report := tc.gf.fetchObject(context.Background(), j)

	if !report.success {
		t.Errorf("report.success got false, want true")
	}
	if report.err != nil {
		t.Errorf("report.err got %v, want nil", report.err)
	}

	if len(report.attempts) != 2 {
		t.Fatalf("len(report.attempts) got %d, want 2", len(report.attempts))
	}
	first := report.attempts[0]
	if first.err == nil || !strings.Contains(first.err.Error(), errChmod.Error()) {
		t.Errorf("attempt.err got %v, want Contains(%v)", first.err, errChmod)
	}

	got, err := ioutil.ReadFile(report.finalname)
	if err != nil {
		t.Fatalf("ReadFile(%v) got %v, want nil", report.finalname, err)
	}
	if !bytes.Equal(got, sfile1Contents) {
		t.Fatalf("ReadFile(%v) got %v, want %v", report.finalname, got, sfile1Contents)
	}
}

func TestDoWork(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()

	files := []string{sfile1, sfile2, sfile3}
	sort.Strings(files)

	// Add n jobs
	todo := make(chan job, len(files))
	results := make(chan jobReport, len(files))
	for i, file := range files {
		todo <- job{bucket: successBucket, object: file, filename: fmt.Sprintf("sfile-%d", i)}
	}

	// Process the jobs
	go tc.gf.doWork(context.Background(), todo, results)

	// Get n reports
	var gotFiles []string
	for range files {
		report := <-results
		if report.err != nil {
			t.Errorf("file %q: report.err got %v, want nil", report.job.filename, report.err)
		}
		if _, err := os.Stat(report.finalname); os.IsNotExist(err) {
			t.Errorf("file %q: does not exist, but it should exist", report.finalname)
		}
		gotFiles = append(gotFiles, report.job.object)
	}

	// Ensure there is nothing more in the results channel
	select {
	case report, ok := <-results:
		if ok {
			t.Errorf("unexpected report found on channel: %v", report)
		} else {
			close(todo)
		}
	default:
	}
	close(results)

	sort.Strings(gotFiles)
	if !reflect.DeepEqual(gotFiles, files) {
		t.Fatalf("processJobs files got %v, want %v", gotFiles, files)
	}
}

func TestProcessJobs(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()
	tc.os.errorsCreate = 1 // Provoke one retry

	jobs := []job{
		{bucket: successBucket, object: sfile1, filename: "sfile1"},
		{bucket: successBucket, object: sfile2, filename: "sfile2"},
		{bucket: successBucket, object: sfile3, filename: "sfile3"},
	}

	stats := tc.gf.processJobs(context.Background(), jobs)

	if !stats.success {
		t.Errorf("processJobs() stats.success got false, want true")
	}
	if len(stats.errs) != 0 {
		t.Errorf("processJobs() stats.errs got %v, want {}", stats.errs)
	}
	if stats.files != len(jobs) {
		t.Errorf("processJobs stats.files got %d, want %d", stats.files, len(jobs))
	}

	wantSize := len(sfile1Contents) + len(sfile2Contents) + len(sfile3Contents)
	if int(stats.size) != wantSize {
		t.Errorf("processJobs() stats.size got %d, want %d", stats.size, wantSize)
	}
	if stats.retries != 1 {
		t.Errorf("processJobs() stats.retries got %d, want 1", stats.retries)
	}
}

func TestFetchFromManifestSuccess(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()

	tc.gf.Bucket = successBucket
	tc.gf.Object = goodManifest

	err := tc.gf.fetchFromManifest(context.Background())
	if err != nil {
		t.Errorf("fetchFromManifest() got %v, want nil", err)
	}

	// Check that enough files are present
	infos, err := ioutil.ReadDir(tc.gf.DestDir)
	if err != nil {
		t.Fatalf("ReadDir(%v) err = %v, want nil", tc.gf.DestDir, err)
	}
	if len(infos) != 4 { // 3 files in the manifest + the manifest itself
		t.Errorf("ReadDir(%v) len(fileinfos)=%v, want 4", tc.gf.DestDir, len(infos))
	}
}

func TestFetchFromManifestManifestGCSFetchFailed(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()

	tc.gf.Bucket = errorBucket
	tc.gf.Object = errorManifest

	err := tc.gf.fetchFromManifest(context.Background())
	if err == nil || !strings.Contains(err.Error(), errGCSRead.Error()) {
		t.Errorf("fetchFromManifest() err=%v, want contains %v", err, errGCSRead)
	}
}

func TestFetchFromManifestManifestJSONDeserializtionFailed(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()

	tc.gf.Bucket = successBucket
	tc.gf.Object = malformedManifest

	wantErrStr := "decoding JSON from manifest file"
	err := tc.gf.fetchFromManifest(context.Background())
	if err == nil || !strings.Contains(err.Error(), wantErrStr) {
		t.Errorf("fetchFromManifest() err=%v, want contains %q", err, wantErrStr)
	}
}

func TestFetchFromManifestManifestFileReadFailed(t *testing.T) {
	tc, teardown := buildManifestTestContext(t)
	defer teardown()
	tc.os.errorsOpen = 1 // Error returned when trying to open the downloaded manifest file

	err := tc.gf.fetchFromManifest(context.Background())
	if err == nil || !strings.Contains(err.Error(), errOpen.Error()) {
		t.Errorf("fetchFromManifest() err=%v, want contains %v", err, errOpen)
	}
}

func TestTimeout(t *testing.T) {
	tests := []struct {
		filename string
		retrynum int
		want     time.Duration
	}{
		{"source.js", 0, sourceTimeout[0]},
		{"source.js", 1, sourceTimeout[1]},
		{"source.js", 2, defaultTimeout},
		{"not-source.mpg", 0, notSourceTimeout[0]},
		{"not-source.mpg", 1, notSourceTimeout[1]},
		{"not-source.mpg", 2, defaultTimeout},
		{"no-extension", 0, notSourceTimeout[0]},
		{"no-extension", 1, notSourceTimeout[1]},
		{"no-extension", 2, defaultTimeout},
	}
	tc, teardown := buildManifestTestContext(t)
	defer teardown()
	for _, test := range tests {
		got := tc.gf.timeout(test.filename, test.retrynum)
		if got != test.want {
			t.Errorf("getTimeout(%v, %v) got %v, want %v", test.filename, test.retrynum, got, test.want)
		}
	}
}

func TestParseBucketObject(t *testing.T) {
	tests := []struct {
		filename string
		bucket   string
		object   string
		err      error
	}{{
		"https://storage.googleapis.com/staging.appid.appspot.com/abc123",
		"staging.appid.appspot.com",
		"abc123",
		nil,
	}, {
		"https://storage.googleapis.com/some-bucket.google.com.a.appspot.com/some/path/to/file",
		"some-bucket.google.com.a.appspot.com",
		"some/path/to/file",
		nil,
	}, {
		"https://storage.googleapis.com/some-bucket/abc123",
		"some-bucket",
		"abc123",
		nil,
	}, {
		"https://storage.googleapis.com/too-short",
		"",
		"",
		errNonNil,
	}, {
		"https://incorrect-domain.com/some-bucket.google.com.a.appspot.com/some/path",
		"",
		"",
		errNonNil,
	}, {
		"gs://my-bucket/manifest-20171004T175409.json",
		"my-bucket",
		"manifest-20171004T175409.json",
		nil,
	}, {
		"gs://staging.appid.appspot.com/abc123",
		"staging.appid.appspot.com",
		"abc123",
		nil,
	}, {
		"gs://some-bucket.google.com.a.appspot.com/some/path/to/file",
		"some-bucket.google.com.a.appspot.com",
		"some/path/to/file",
		nil,
	}, {
		"gs://some-bucket/abc123",
		"some-bucket",
		"abc123",
		nil,
	}, {
		"http://storage.googleapis.com/my-bucket/test-memchache/server.js",
		"my-bucket",
		"test-memchache/server.js",
		nil,
	}, {
		"gs://too-short",
		"",
		"",
		errNonNil,
	}, {
		"some-bucket/some/path",
		"",
		"",
		errNonNil,
	}}
	for _, test := range tests {
		bucket, object, _, err := ParseBucketObject(test.filename)
		if test.err == nil {
			if bucket != test.bucket || object != test.object {
				t.Errorf("parseBucketObject(%q) = (%q, %q); want (%q, %q)", test.filename, bucket, object, test.bucket, test.object)
			}
		} else {
			if err == nil {
				t.Errorf("parseBucketObject(%q) = (%q, %q), want err", test.filename, bucket, object)
			}
		}
	}
}
