package podmanProv

import (
	"context"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/kigland/OpenHPC/lib/consts"
)

type PodmanProvider struct {
	conn context.Context
}

func NewPodmanSession() (*PodmanProvider, error) {
	conn, err := bindings.NewConnection(context.Background(), consts.PODMAN_UNIX_SOCKET)
	if err != nil {
		return nil, err
	}
	return &PodmanProvider{conn: conn}, nil
}

func (p *PodmanProvider) Cli() context.Context {
	return p.conn
}
