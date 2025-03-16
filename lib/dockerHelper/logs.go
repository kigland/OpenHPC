package dockerHelper

import (
	"bytes"
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
)

func (d *DockerHelper) GetLogs(containerID string, timestamps bool) (string, error) {
	v, err := d.GetLogsRaw(containerID, timestamps)
	if err != nil {
		return "", err
	}
	defer v.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, v)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (d *DockerHelper) GetLogsRaw(containerID string, timestamps bool) (io.ReadCloser, error) {
	v, err := d.cli.ContainerLogs(context.Background(), containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Timestamps: timestamps,
		Details:    true,
	})
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (d *DockerHelper) GetLogsStream(containerID string, timestamps bool) (io.ReadCloser, error) {
	return d.cli.ContainerLogs(context.Background(), containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: timestamps,
		Details:    true,
	})
}
