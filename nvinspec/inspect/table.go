package inspect

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

func WriteTable(w io.Writer, findings []Finding, fullID bool) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "no active nvidia-smi compute processes found")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "GPU\tPID\tPROCESS\tMEM(MiB)\tCONTAINER\tNAME\tIMAGE\tSOURCE")
	for _, finding := range findings {
		containerID := finding.ContainerID
		if containerID == "" {
			containerID = "-"
		} else if !fullID {
			containerID = ShortContainerID(containerID)
		}

		mem := "-"
		if finding.UsedMemoryMiB >= 0 {
			mem = fmt.Sprintf("%d", finding.UsedMemoryMiB)
		}

		fmt.Fprintf(
			tw,
			"%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
			shortGPUUUID(finding.GPUUUID),
			finding.PID,
			emptyDash(finding.ProcessName),
			mem,
			containerID,
			emptyDash(finding.ContainerName),
			emptyDash(finding.ContainerImage),
			finding.Source,
		)
	}
	_ = tw.Flush()
}

func shortGPUUUID(uuid string) string {
	const prefix = "GPU-"
	uuid = strings.TrimSpace(uuid)
	if len(uuid) <= 12 {
		return emptyDash(uuid)
	}
	if strings.HasPrefix(uuid, prefix) && len(uuid) > len(prefix)+8 {
		return uuid[:len(prefix)+8]
	}
	return uuid[:12]
}

func emptyDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}
