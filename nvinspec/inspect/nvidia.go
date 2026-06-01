package inspect

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

func queryGPUProcesses(ctx context.Context, runner Runner, opts Options) ([]GPUProcess, error) {
	out, err := runner.Run(
		ctx,
		opts.NvidiaSMI,
		"--query-compute-apps=gpu_uuid,pid,process_name,used_memory",
		"--format=csv,noheader,nounits",
	)
	if err != nil {
		return nil, fmt.Errorf("query nvidia-smi compute processes: %w", err)
	}
	return ParseNvidiaSmiCSV(out)
}

func ParseNvidiaSmiCSV(data []byte) ([]GPUProcess, error) {
	text := strings.TrimSpace(string(data))
	if text == "" || strings.Contains(strings.ToLower(text), "no running processes") {
		return nil, nil
	}

	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse nvidia-smi csv: %w", err)
	}

	processes := make([]GPUProcess, 0, len(records))
	for line, record := range records {
		if len(record) == 0 {
			continue
		}
		if len(record) == 1 && strings.TrimSpace(record[0]) == "" {
			continue
		}
		if len(record) != 4 {
			return nil, fmt.Errorf("parse nvidia-smi csv line %d: expected 4 fields, got %d", line+1, len(record))
		}

		pid, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			return nil, fmt.Errorf("parse nvidia-smi csv line %d pid: %w", line+1, err)
		}

		processes = append(processes, GPUProcess{
			GPUUUID:       strings.TrimSpace(record[0]),
			PID:           pid,
			ProcessName:   strings.TrimSpace(record[2]),
			UsedMemoryMiB: parseMemoryMiB(record[3]),
		})
	}
	return processes, nil
}

func parseMemoryMiB(raw string) int {
	value := strings.TrimSpace(raw)
	value = strings.TrimSuffix(value, "MiB")
	value = strings.TrimSpace(value)
	if value == "" || strings.Contains(value, "Not Supported") || value == "[N/A]" || value == "N/A" {
		return -1
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return -1
	}
	return n
}
