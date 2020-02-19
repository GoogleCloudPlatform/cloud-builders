package services

import (
	"context"
	"fmt"
	"os/exec"
)

const GsutilExec = "gsutil"

// Gsutil implements the GcsService interface
type Gsutil struct {
	printCommands bool
}

// NewGsutil returns a new Gsutil object
func NewGsutil(ctx context.Context, printCommands bool) (*Gsutil, error) {
	if _, err := exec.LookPath(GsutilExec); err != nil {
		return nil, err
	}
	return &Gsutil{printCommands: printCommands}, nil
}

// Copy calls `gsutil cp -r <source_url> <destination_url>
func (g *Gsutil) Copy(ctx context.Context, src, dst string, recursive bool) error {
	args := []string{"cp", "-r", src, dst}
	// remove the -r flag
	if !recursive {
		args = append(args[:1], args[2:]...)
	}
	if _, err := runCommand(ctx, g.printCommands, GsutilExec, args...); err != nil {
		return fmt.Errorf("copy file(s) with %s failed: %v", GsutilExec, err.Error())
	}
	return nil
}
