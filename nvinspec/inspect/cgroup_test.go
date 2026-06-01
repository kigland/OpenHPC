package inspect

import "testing"

func TestExtractPodmanContainerIDFromCgroup(t *testing.T) {
	id := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	cgroup := "0::/user.slice/user-1000.slice/user@1000.service/user.slice/libpod-" + id + ".scope\n"

	got := ExtractPodmanContainerIDFromCgroup(cgroup)
	if got != id {
		t.Fatalf("ExtractPodmanContainerIDFromCgroup() = %q, want %q", got, id)
	}
}

func TestExtractPodmanContainerIDFromCgroupCrun(t *testing.T) {
	id := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	cgroup := "0::/machine.slice/crun-" + id + ".scope\n"

	got := ExtractPodmanContainerIDFromCgroup(cgroup)
	if got != id {
		t.Fatalf("ExtractPodmanContainerIDFromCgroup() = %q, want %q", got, id)
	}
}

func TestExtractPodmanContainerIDFromCgroupIgnoresDocker(t *testing.T) {
	id := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	cgroup := "0::/docker/" + id + "\n"

	got := ExtractPodmanContainerIDFromCgroup(cgroup)
	if got != "" {
		t.Fatalf("ExtractPodmanContainerIDFromCgroup() = %q, want empty", got)
	}
}
