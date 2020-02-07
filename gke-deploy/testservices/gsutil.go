package testservices

import (
	"context"
	"fmt"
)

// GCSResponse simulates errors while calling GcsService.
//type GCSResponse interface {
//	Error() error
//}

type TestGcsService struct {
	CopyResponse map[string]func(src, dst string, recursive bool) error
}

func (s *TestGcsService) Copy(ctx context.Context, src, dst string, recursive bool) error {
	res, ok := s.CopyResponse[src]
	if !ok {
		res, ok = s.CopyResponse[dst]
		if !ok {
			panic(fmt.Sprintf("no response for source %q", src))
		}
	}
	return res(src, dst, recursive)
}
