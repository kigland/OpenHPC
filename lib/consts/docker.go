package consts

import "strconv"

const (
	DOCKER_UNIX_SOCKET = "unix:///var/run/docker.sock"
	PODMAN_UNIX_SOCKET = "unix:///run/podman/podman.sock"
)

const (
	RDS_PATH = "/data/rds"
)

const (
	IDENTIFIER = "KHS"
)

const LOW_PORT = 40_000
const HIGH_PORT = 40_499
const SSH_SHIFT = 500
const CONTAINER_HOST = "127.0.0.2"
const SSH_BIND_HOST = "127.0.0.3"

func BASE_URL(port int) string {
	return "/ohpc/" + strconv.Itoa(port)
}
