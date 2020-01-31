package services

import (
	"context"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Remote implements the RemoteService interface.
type Remote struct{}

// NewRemote returns a new Remote object.
func NewRemote(ctx context.Context) (*Remote, error) {
	return &Remote{}, nil
}

// Image gets a remote image from a reference.
func (*Remote) Image(ref name.Reference) (v1.Image, error) {
	return remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
}
