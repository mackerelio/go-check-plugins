//go:build docker
// +build docker

package checkldap

import (
	"testing"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/mackerelio/checkers"
)

var imageRepository = "osixia/openldap"
var imageTag = "1.2.0"

func TestLDAP(t *testing.T) {
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: imageRepository + ":" + imageTag,
			ExposedPorts: map[docker.Port]struct{}{
				"389/tcp": {},
			},
		},
		HostConfig: &docker.HostConfig{
			PortBindings: map[docker.Port][]docker.PortBinding{
				"389/tcp": {
					{
						HostIP:   "",
						HostPort: "389",
					},
				},
			},
		},
	}

	cli, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = cli.InspectImage(opts.Config.Image)
	if err != nil {
		err = cli.PullImage(docker.PullImageOptions{
			Repository: imageRepository,
			Tag:        imageTag,
		}, docker.AuthConfiguration{})
		if err != nil {
			t.Fatal(err)
		}
	}

	container, err := cli.CreateContainer(opts)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := cli.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		})
		if err != nil {
			t.Fatal(err)
		}
	}()

	if err := cli.StartContainer(container.ID, &docker.HostConfig{}); err != nil {
		t.Fatal(err)
	}

	// Wait for initialization of LDAP server
	time.Sleep(3 * time.Second)

	args := []string{"-H", "localhost", "-c", "1", "-w", "0.5", "-b", "dc=example,dc=org", "-D", "cn=admin,dc=example,dc=org", "-P", "admin"}

	ckr := run(args)

	if ckr.Status != checkers.OK {
		t.Errorf("ckr.Status should be OK: %s", ckr)
	}
}
