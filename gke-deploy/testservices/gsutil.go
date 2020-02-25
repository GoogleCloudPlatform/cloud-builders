package testservices

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const tmpPermissionForDirectory = os.FileMode(0755)

type TestGcsService struct {
	CopyResponse map[string]func(src, dst string) error
}

func (s *TestGcsService) Copy(ctx context.Context, src, dst string, recursive bool) error {
	res, ok := s.CopyResponse[src]
	if !ok {
		res, ok = s.CopyResponse[dst]
		if !ok {
			panic(fmt.Sprintf("no response for source %q", src))
		}
	}
	return res(src, dst)
}

func copyFile(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

func copy(src, dest string, info os.FileInfo) error {
	if info.IsDir() {
		return copyDir(src, dest, info)
	}
	if !strings.HasSuffix(dest, "yaml") && !strings.HasSuffix(dest, "yml") {
		dest = filepath.Join(dest, info.Name())
	}
	return copyFile(src, dest)
}

func copyDir(srcdir, destdir string, info os.FileInfo) error {
	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}
	for _, content := range contents {
		cs := filepath.Join(srcdir, content.Name())

		cd := destdir
		if content.IsDir() {
			cd = filepath.Join(destdir, content.Name())
		}
		// Make dest dir with 0755 so that everything writable.
		if err := os.MkdirAll(cd, tmpPermissionForDirectory); err != nil {
			return err
		}
		if err := copy(cs, cd, content); err != nil {
			return err
		}
	}
	return nil
}

// Copy simulates the gsutil copy.
func Copy(src, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return copy(src, dest, info)
}
