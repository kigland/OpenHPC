package dockerProv

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
)

func (d *DockerHelper) Pull(imageName string) (err error) {
	cli := d.cli
	out, err := cli.ImagePull(context.Background(), imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return err
}
