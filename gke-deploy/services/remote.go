package services

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"golang.org/x/oauth2/google"
)

// Remote implements the RemoteService interface.
type Remote struct{}

// NewRemote returns a new Remote object.
func NewRemote(ctx context.Context) (*Remote, error) {
	return &Remote{}, nil
}

// Image gets a remote image from a reference.
func (*Remote) Image(ctx context.Context, ref name.Reference) (v1.Image, error) {
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/cloud-platform")
		if err != nil {
			fmt.Printf("Error fetching digest: %v\n", err)
			return nil, err
		}
		return remote.Image(ref, remote.WithTransport(client.Transport))
	}
	return img, nil
}
