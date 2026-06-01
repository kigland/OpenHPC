package inspect

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type containerCatalog struct {
	containers []ContainerInfo
}

func loadPodmanContainers(ctx context.Context, runner Runner, opts Options) (containerCatalog, error) {
	name, args := podmanCommand(opts, "ps", "--all", "--no-trunc", "--format", "json")
	out, err := runner.Run(ctx, name, args...)
	if err != nil {
		return containerCatalog{}, fmt.Errorf("podman ps: %w", err)
	}
	return ParsePodmanPSJSON(out)
}

func ParsePodmanPSJSON(data []byte) (containerCatalog, error) {
	data = []byte(strings.TrimSpace(string(data)))
	if len(data) == 0 || string(data) == "null" {
		return containerCatalog{}, nil
	}

	var rows []map[string]any
	if err := json.Unmarshal(data, &rows); err != nil {
		return containerCatalog{}, fmt.Errorf("parse podman ps json: %w", err)
	}

	catalog := containerCatalog{
		containers: make([]ContainerInfo, 0, len(rows)),
	}
	for _, row := range rows {
		info := ContainerInfo{
			ID:    firstString(row, "Id", "ID", "ContainerID"),
			Name:  firstName(row, "Names", "Name"),
			Image: firstString(row, "Image", "ImageName"),
		}
		if info.ID == "" {
			continue
		}
		catalog.containers = append(catalog.containers, info)
	}
	return catalog, nil
}

func (c containerCatalog) find(id string) (ContainerInfo, bool) {
	if id == "" {
		return ContainerInfo{}, false
	}
	id = strings.ToLower(id)
	for _, container := range c.containers {
		cID := strings.ToLower(container.ID)
		if cID == id || strings.HasPrefix(cID, id) || strings.HasPrefix(id, cID) {
			return container, true
		}
	}
	return ContainerInfo{ID: id}, false
}

func firstString(row map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := row[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case string:
			if typed != "" {
				return strings.TrimPrefix(typed, "/")
			}
		case fmt.Stringer:
			if typed.String() != "" {
				return strings.TrimPrefix(typed.String(), "/")
			}
		}
	}
	return ""
}

func firstName(row map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := row[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case string:
			return strings.TrimPrefix(typed, "/")
		case []any:
			for _, item := range typed {
				if name, ok := item.(string); ok && name != "" {
					return strings.TrimPrefix(name, "/")
				}
			}
		case []string:
			if len(typed) > 0 {
				return strings.TrimPrefix(typed[0], "/")
			}
		}
	}
	return ""
}

func buildPodmanPIDMap(ctx context.Context, runner Runner, opts Options, catalog containerCatalog) (map[int]ContainerInfo, []string) {
	out := map[int]ContainerInfo{}
	warnings := []string{}

	for _, container := range catalog.containers {
		if container.ID == "" {
			continue
		}

		data, err := runPodmanTop(ctx, runner, opts, container.ID, "args")
		if err != nil {
			data, err = runPodmanTop(ctx, runner, opts, container.ID, "comm")
		}
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("podman top %s: %v", ShortContainerID(container.ID), err))
			continue
		}

		for _, pid := range ParsePodmanTopHPIDs(data) {
			out[pid] = container
		}
	}
	return out, warnings
}

func runPodmanTop(ctx context.Context, runner Runner, opts Options, id string, commandColumn string) ([]byte, error) {
	name, args := podmanCommand(opts, "top", id, "hpid", "pid", commandColumn)
	return runner.Run(ctx, name, args...)
}

func ParsePodmanTopHPIDs(data []byte) []int {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) <= 1 {
		return nil
	}

	pids := make([]int, 0, len(lines)-1)
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		pids = append(pids, pid)
	}
	return pids
}
