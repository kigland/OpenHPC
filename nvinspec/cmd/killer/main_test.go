package main

import (
	"testing"

	"github.com/kigland/OpenHPC/nvinspec/inspect"
)

func TestPlanVictimsKillsLargestUntilStrictlyBelowLimit(t *testing.T) {
	findings := []inspect.Finding{
		{PID: 1, ContainerID: "b3d495dd4772aaaaaaaa", UsedMemoryMiB: 4096},
		{PID: 2, ContainerID: "b3d495dd4772aaaaaaaa", UsedMemoryMiB: 12288},
		{PID: 3, ContainerID: "b3d495dd4772aaaaaaaa", UsedMemoryMiB: 8192},
		{PID: 4, ContainerID: "another-container", UsedMemoryMiB: 30000},
	}

	total, victims := planVictims(findings, "b3d495dd4772", 20*1024)
	if total != 24*1024 {
		t.Fatalf("total = %d, want %d", total, 24*1024)
	}
	if len(victims) != 1 || victims[0].PID != 2 {
		t.Fatalf("victims = %+v, want PID 2", victims)
	}
}

func TestPlanVictimsDoesNothingAtLimit(t *testing.T) {
	findings := []inspect.Finding{
		{PID: 1, ContainerID: defaultContainerID, UsedMemoryMiB: defaultLimitMiB},
	}

	total, victims := planVictims(findings, defaultContainerID, defaultLimitMiB)
	if total != defaultLimitMiB || len(victims) != 0 {
		t.Fatalf("total = %d, victims = %+v", total, victims)
	}
}

func TestPlanVictimsKillsPastExactLimitAfterTrigger(t *testing.T) {
	findings := make([]inspect.Finding, 0, 21)
	for pid := 1; pid <= 21; pid++ {
		findings = append(findings, inspect.Finding{
			PID:           pid,
			ContainerID:   defaultContainerID,
			UsedMemoryMiB: 1024,
		})
	}

	total, victims := planVictims(findings, defaultContainerID, defaultLimitMiB)
	if total != 21*1024 {
		t.Fatalf("total = %d, want %d", total, 21*1024)
	}
	if len(victims) != 2 {
		t.Fatalf("len(victims) = %d, want 2 so projected usage is strictly below the limit", len(victims))
	}
}

func TestSameContainerID(t *testing.T) {
	full := "b3d495dd4772abcdef0123456789"
	if !sameContainerID(full, defaultContainerID) {
		t.Fatal("full ID should match configured prefix")
	}
	if sameContainerID("", defaultContainerID) {
		t.Fatal("empty ID must not match")
	}
	if sameContainerID("abcdef", defaultContainerID) {
		t.Fatal("different ID must not match")
	}
}
