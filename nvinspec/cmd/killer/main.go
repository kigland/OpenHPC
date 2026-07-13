package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kigland/OpenHPC/nvinspec/inspect"
)

const (
	defaultContainerID = "b3d495dd4772"
	defaultLimitMiB    = 20 * 1024
)

type config struct {
	inspectOptions inspect.Options
	containerID    string
	limitMiB       int
	interval       time.Duration
	once           bool
}

type killFunc func(context.Context, int) error

func main() {
	cfg := parseFlags()
	if err := validateConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "nvinspec-killer: %v\n", err)
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		if err := enforceLimit(ctx, cfg, kill9); err != nil {
			if ctx.Err() != nil {
				return
			}
			fmt.Fprintf(os.Stderr, "nvinspec-killer: %v\n", err)
			if cfg.once {
				os.Exit(1)
			}
		}

		if cfg.once || !wait(ctx, cfg.interval) {
			return
		}
	}
}

func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.inspectOptions.NvidiaSMI, "nvidia-smi", "nvidia-smi", "path to nvidia-smi")
	flag.StringVar(&cfg.inspectOptions.Podman, "podman", "podman", "path to podman")
	flag.StringVar(&cfg.inspectOptions.ProcRoot, "proc-root", "/proc", "procfs root")
	flag.BoolVar(&cfg.inspectOptions.UseSudo, "sudo", false, "run podman through sudo")
	flag.DurationVar(&cfg.inspectOptions.Timeout, "timeout", 8*time.Second, "timeout for each inspection")
	flag.StringVar(&cfg.containerID, "container", defaultContainerID, "target Podman container ID or ID prefix")
	flag.IntVar(&cfg.limitMiB, "limit-mib", defaultLimitMiB, "GPU memory limit in MiB")
	flag.DurationVar(&cfg.interval, "interval", 10*time.Second, "monitoring interval")
	flag.BoolVar(&cfg.once, "once", false, "enforce the limit once and exit")
	flag.Parse()
	return cfg
}

func validateConfig(cfg config) error {
	if strings.TrimSpace(cfg.containerID) == "" {
		return fmt.Errorf("container ID must not be empty")
	}
	if cfg.limitMiB <= 0 {
		return fmt.Errorf("limit-mib must be greater than zero")
	}
	if cfg.interval <= 0 {
		return fmt.Errorf("interval must be greater than zero")
	}
	if cfg.inspectOptions.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}
	return nil
}

func enforceLimit(ctx context.Context, cfg config, kill killFunc) error {
	for {
		scanCtx, cancel := context.WithTimeout(ctx, cfg.inspectOptions.Timeout)
		report, err := inspect.Run(scanCtx, cfg.inspectOptions)
		cancel()
		if err != nil {
			return fmt.Errorf("inspect GPU processes: %w", err)
		}
		for _, warning := range report.Warnings {
			fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
		}

		totalMiB, victims := planVictims(report.Findings, cfg.containerID, cfg.limitMiB)
		if totalMiB <= cfg.limitMiB {
			fmt.Printf("container %s GPU memory: %d MiB (limit %d MiB)\n", cfg.containerID, totalMiB, cfg.limitMiB)
			return nil
		}
		if len(victims) == 0 {
			return fmt.Errorf("container %s uses %d MiB, but no process with known positive GPU memory can be killed", cfg.containerID, totalMiB)
		}

		remainingMiB := totalMiB
		for _, victim := range victims {
			// The PID came from a fresh inspection that also verified its container.
			fmt.Printf(
				"container %s projected GPU memory %d MiB (target < %d MiB); kill -9 PID %d (%s, %d MiB)\n",
				cfg.containerID,
				remainingMiB,
				cfg.limitMiB,
				victim.PID,
				victim.ProcessName,
				victim.UsedMemoryMiB,
			)
			if err := kill(ctx, victim.PID); err != nil {
				return err
			}
			remainingMiB -= victim.UsedMemoryMiB
		}

		// Give the NVIDIA driver a moment to remove killed processes, then
		// inspect again instead of trusting the projected memory total.
		if !wait(ctx, 500*time.Millisecond) {
			return ctx.Err()
		}
	}
}

func planVictims(findings []inspect.Finding, containerID string, limitMiB int) (int, []inspect.Finding) {
	totalMiB := 0
	candidates := make([]inspect.Finding, 0)
	for _, finding := range findings {
		if !sameContainerID(finding.ContainerID, containerID) || finding.UsedMemoryMiB < 0 {
			continue
		}
		totalMiB += finding.UsedMemoryMiB
		if finding.UsedMemoryMiB > 0 {
			candidates = append(candidates, finding)
		}
	}

	if totalMiB <= limitMiB {
		return totalMiB, nil
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].UsedMemoryMiB > candidates[j].UsedMemoryMiB
	})

	remainingMiB := totalMiB
	victims := make([]inspect.Finding, 0, len(candidates))
	for _, candidate := range candidates {
		if remainingMiB < limitMiB {
			break
		}
		victims = append(victims, candidate)
		remainingMiB -= candidate.UsedMemoryMiB
	}
	return totalMiB, victims
}

func sameContainerID(left, right string) bool {
	left = strings.ToLower(strings.TrimSpace(left))
	right = strings.ToLower(strings.TrimSpace(right))
	return left != "" && right != "" && (left == right || strings.HasPrefix(left, right) || strings.HasPrefix(right, left))
}

func kill9(ctx context.Context, pid int) error {
	cmd := exec.CommandContext(ctx, "kill", "-9", strconv.Itoa(pid))
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	message := strings.TrimSpace(string(out))
	if message == "" {
		return fmt.Errorf("kill -9 PID %d: %w", pid, err)
	}
	return fmt.Errorf("kill -9 PID %d: %w: %s", pid, err, message)
}

func wait(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
