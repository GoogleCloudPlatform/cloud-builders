// Package image contains related to container images.
package image

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
	"github.com/google/go-containerregistry/pkg/name"
)

// Name gets an image's name from a Reference.
// e.g., If the string representation of the Reference is "gcr.io/my-project/my-image:1.0.0", this
// returns "gcr.io/my-project/my-image".
func Name(ref name.Reference) string {
	return fmt.Sprintf("%s/%s", ref.Context().RegistryStr(), ref.Context().RepositoryStr())
}

// ResolveDigest gets an image's corresponding digest.
func ResolveDigest(ctx context.Context, ref name.Reference, rs services.RemoteService) (string, error) {
	im, err := rs.Image(ref)
	if err != nil {
		return "", fmt.Errorf("failed to get remote image reference: %v", err)
	}
	digest, err := im.Digest()
	if err != nil {
		return "", fmt.Errorf("failed to get image digest: %v", err)
	}
	return fmt.Sprintf("%s:%s", digest.Algorithm, digest.Hex), nil
}
