package inspect

import (
	"context"
)

func Run(ctx context.Context, opts Options) (Report, error) {
	opts.setDefaults()
	return RunWithRunner(ctx, opts, OSRunner{})
}

func RunWithRunner(ctx context.Context, opts Options, runner Runner) (Report, error) {
	opts.setDefaults()

	processes, err := queryGPUProcesses(ctx, runner, opts)
	if err != nil {
		return Report{}, err
	}

	report := Report{
		Findings: make([]Finding, 0, len(processes)),
	}

	catalog, err := loadPodmanContainers(ctx, runner, opts)
	if err != nil {
		report.Warnings = append(report.Warnings, err.Error())
	}

	for _, process := range processes {
		finding := Finding{
			GPUUUID:       process.GPUUUID,
			PID:           process.PID,
			ProcessName:   process.ProcessName,
			UsedMemoryMiB: process.UsedMemoryMiB,
			Source:        "unknown",
		}

		containerID, reason := containerIDFromProcCgroup(opts.ProcRoot, process.PID)
		if containerID != "" {
			info, ok := catalog.find(containerID)
			finding.ContainerID = info.ID
			if finding.ContainerID == "" {
				finding.ContainerID = containerID
			}
			finding.ContainerName = info.Name
			finding.ContainerImage = info.Image
			finding.Source = "cgroup"
			if !ok {
				finding.Reason = "container id found in cgroup; podman metadata unavailable"
			}
		} else {
			finding.Reason = reason
		}

		report.Findings = append(report.Findings, finding)
	}

	fillFromPodmanTop(ctx, runner, opts, catalog, &report)
	return report, nil
}

func fillFromPodmanTop(ctx context.Context, runner Runner, opts Options, catalog containerCatalog, report *Report) {
	if len(catalog.containers) == 0 {
		return
	}

	needsFallback := false
	for _, finding := range report.Findings {
		if finding.ContainerID == "" {
			needsFallback = true
			break
		}
	}
	if !needsFallback {
		return
	}

	pidMap, warnings := buildPodmanPIDMap(ctx, runner, opts, catalog)
	report.Warnings = append(report.Warnings, warnings...)
	if len(pidMap) == 0 {
		return
	}

	for i := range report.Findings {
		if report.Findings[i].ContainerID != "" {
			continue
		}
		container, ok := pidMap[report.Findings[i].PID]
		if !ok {
			continue
		}
		report.Findings[i].ContainerID = container.ID
		report.Findings[i].ContainerName = container.Name
		report.Findings[i].ContainerImage = container.Image
		report.Findings[i].Source = "podman-top"
		report.Findings[i].Reason = ""
	}
}
