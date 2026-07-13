package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kigland/OpenHPC/nvinspec/inspect"
)

func main() {
	var opts inspect.Options
	var jsonOutput bool
	var fullID bool
	var containersOnly bool

	flag.StringVar(&opts.NvidiaSMI, "nvidia-smi", "nvidia-smi", "path to nvidia-smi")
	flag.StringVar(&opts.Podman, "podman", "podman", "path to podman")
	flag.StringVar(&opts.ProcRoot, "proc-root", "/proc", "procfs root")
	flag.BoolVar(&opts.UseSudo, "sudo", false, "run podman through sudo")
	flag.BoolVar(&jsonOutput, "json", false, "write JSON output")
	flag.BoolVar(&fullID, "full-id", false, "show full container IDs in table output")
	flag.BoolVar(&containersOnly, "containers-only", false, "hide GPU processes that are not mapped to a Podman container")
	flag.DurationVar(&opts.Timeout, "timeout", 8*time.Second, "command timeout")
	flag.Usage = usage
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	report, err := inspect.Run(ctx, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nvinspec: %v\n", err)
		os.Exit(1)
	}

	if containersOnly {
		report.Findings = filterContainersOnly(report.Findings)
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			fmt.Fprintf(os.Stderr, "nvinspec: write json: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for _, warning := range report.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}
	inspect.WriteTable(os.Stdout, report.Findings, fullID)
}

func filterContainersOnly(findings []inspect.Finding) []inspect.Finding {
	out := findings[:0]
	for _, finding := range findings {
		if finding.ContainerID != "" {
			out = append(out, finding)
		}
	}
	return out
}

func usage() {
	h := `
Usage:
  nvinspec [flags]

Detect Podman containers for GPU processes reported by nvidia-smi.

Flags:
`
	fmt.Fprint(flag.CommandLine.Output(), strings.TrimLeft(h, "\n"))
	flag.PrintDefaults()
}
