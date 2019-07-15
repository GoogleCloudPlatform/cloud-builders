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
// Package testservices contains fake services used for unit tests.
package testservices

import (
	"context"
	"fmt"
	"os"
)

// TestOS implements the OSService interface.
type TestOS struct {
	StatResponse      map[string]StatResponse
	ReadDirResponse   map[string]ReadDirResponse
	ReadFileResponse  map[string]ReadFileResponse
	WriteFileResponse map[string]error
	MkdirAllResponse  map[string]error
}

// StatResponse represents a response tuple for a Stat function call.
type StatResponse struct {
	Res os.FileInfo
	Err error
}

// ReadDirResponse represents a response tuple for a ReadDir function call.
type ReadDirResponse struct {
	Res []os.FileInfo
	Err error
}

// ReadDirResponse represents a response tuple for a ReadFile function call.
type ReadFileResponse struct {
	Res []byte
	Err error
}

// Stat gets a file description for a file filename.
func (o *TestOS) Stat(ctx context.Context, filename string) (os.FileInfo, error) {
	resp, ok := o.StatResponse[filename]
	if !ok {
		panic(fmt.Sprintf("Stat has no response for filename %q", filename))
	}
	return resp.Res, resp.Err
}

// ReadDir gets file descriptions for all files contained in a directory dirname.
func (o *TestOS) ReadDir(ctx context.Context, dirname string) ([]os.FileInfo, error) {
	resp, ok := o.ReadDirResponse[dirname]
	if !ok {
		panic(fmt.Sprintf("ReadDir has no response for dirname %q", dirname))
	}
	return resp.Res, resp.Err
}

// ReadFile gets the entire contents of a file filename as bytes.
func (o *TestOS) ReadFile(ctx context.Context, filename string) ([]byte, error) {
	resp, ok := o.ReadFileResponse[filename]
	if !ok {
		panic(fmt.Sprintf("ReadFile has no response for filename %q", filename))
	}
	return resp.Res, resp.Err
}

// WriteFile writes data to a file.
func (o *TestOS) WriteFile(ctx context.Context, filename string, data []byte, perm os.FileMode) error {
	err, ok := o.WriteFileResponse[filename]
	if !ok {
		panic(fmt.Sprintf("WriteFileResponse has no response for filename %q", filename))
	}
	return err
}

// MkdirAll creates a directory dirname, including all parent directories if they do not exist.
func (o *TestOS) MkdirAll(ctx context.Context, dirname string, perm os.FileMode) error {
	err, ok := o.MkdirAllResponse[dirname]
	if !ok {
		panic(fmt.Sprintf("MkdirAllResponse has no response for dirname %q", dirname))
	}
	return err
}

// TestFileInfo implements the os.FileInfo interface.
type TestFileInfo struct {
	// Embed this so we only need to add methods used by testable functions
	os.FileInfo
	BaseName    string
	IsDirectory bool
}

func (fi *TestFileInfo) Name() string {
	return fi.BaseName
}

// IsDir returns true if the file is a directory.
func (fi *TestFileInfo) IsDir() bool {
	return fi.IsDirectory
}
