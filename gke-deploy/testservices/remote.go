package testservices

import (
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
)

// TestRemote implements the RemoteService interface.
type TestRemote struct {
	ImageResp v1.Image
	ImageErr  error
}

// Image gets a remote image from a reference.
func (r *TestRemote) Image(ref name.Reference) (v1.Image, error) {
	return r.ImageResp, r.ImageErr
}

// TestImage simplements the v1.Image interface.
type TestImage struct {
	// Embed this so we only need to add methods used by testable functions
	v1.Image
	Hash v1.Hash
	Err  error
}

// Digest returns the sha256 of this image's manifest.
func (i TestImage) Digest() (v1.Hash, error) {
	return i.Hash, i.Err
}
