package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func runCommandWithStdinRedirection(ctx context.Context, printCommand bool, name, input string, args ...string) (string, error) {
	if printCommand {
		fmt.Printf("\n--------------------------------------------------------------------------------\n")
		fmt.Printf("> Running command\n\n")
		fmt.Printf("   %s %s < %s\n", name, strings.Join(args, " "), input)
		fmt.Printf("\n--------------------------------------------------------------------------------\n\n")
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stderr = os.Stderr
	// Example taken from https://golang.org/src/os/exec/example_test.go
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func runCommand(ctx context.Context, printCommand bool, name string, args ...string) (string, error) {
	if printCommand {
		fmt.Printf("\n--------------------------------------------------------------------------------\n")
		fmt.Printf("> Running command\n\n")
		fmt.Printf("   %s %s\n", name, strings.Join(args, " "))
		fmt.Printf("\n--------------------------------------------------------------------------------\n\n")
	}
	cmd := exec.CommandContext(ctx, name, args...)
	var buf bytes.Buffer
	cmd.Stderr = &buf
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Errorf("%s %s", buf.String(), err.Error())
	}
	return string(out), nil
}
