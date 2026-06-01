package inspect

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

type OSRunner struct{}

func (OSRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err == nil {
		return out, nil
	}
	msg := strings.TrimSpace(stderr.String())
	if msg == "" {
		return out, err
	}
	return out, fmt.Errorf("%w: %s", err, msg)
}

func podmanCommand(opts Options, args ...string) (string, []string) {
	if opts.UseSudo {
		next := append([]string{opts.Podman}, args...)
		return "sudo", next
	}
	return opts.Podman, args
}
