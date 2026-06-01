package inspect

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var podmanCgroupPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)libpod[-_]?([0-9a-f]{12,64})`),
	regexp.MustCompile(`(?i)crun[-_]?([0-9a-f]{12,64})`),
	regexp.MustCompile(`(?i)podman[-_]?([0-9a-f]{12,64})`),
}

var genericContainerIDPattern = regexp.MustCompile(`(?i)(?:^|/)([0-9a-f]{64})(?:/|$|\.scope)`)

func containerIDFromProcCgroup(procRoot string, pid int) (string, string) {
	path := filepath.Join(procRoot, strconv.Itoa(pid), "cgroup")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err.Error()
	}
	id := ExtractPodmanContainerIDFromCgroup(string(data))
	if id == "" {
		return "", "no podman container id in cgroup"
	}
	return id, ""
}

func ExtractPodmanContainerIDFromCgroup(data string) string {
	for _, pattern := range podmanCgroupPatterns {
		if match := pattern.FindStringSubmatch(data); len(match) == 2 {
			return strings.ToLower(match[1])
		}
	}

	lower := strings.ToLower(data)
	if strings.Contains(lower, "libpod") || strings.Contains(lower, "podman") || strings.Contains(lower, "crun") {
		if match := genericContainerIDPattern.FindStringSubmatch(data); len(match) == 2 {
			return strings.ToLower(match[1])
		}
	}
	return ""
}

func ShortContainerID(id string) string {
	if len(id) <= 12 {
		return id
	}
	return id[:12]
}
