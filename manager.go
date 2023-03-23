package inchori_e2e

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type Manager struct {
	pool        *dockertest.Pool
	network     *dockertest.Network
	CurrentNode *dockertest.Resource
}

func NewManager(networkName string) (*Manager, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	network, err := pool.CreateNetwork(networkName)
	if err != nil {
		return nil, err
	}

	return &Manager{
		pool:    pool,
		network: network,
	}, nil
}

func (m *Manager) BuildImage(name, version, dockerFile, contextDir string, args map[string]string) error {
	buildArgs := make([]docker.BuildArg, 0, len(args))
	for k, v := range args {
		bArg := docker.BuildArg{
			Name:  k,
			Value: v,
		}
		buildArgs = append(buildArgs, bArg)
	}
	opts := docker.BuildImageOptions{
		Dockerfile:   dockerFile,
		BuildArgs:    buildArgs,
		NoCache:      true,
		Name:         fmt.Sprintf("%s:%s", name, version),
		OutputStream: io.Discard,
		ErrorStream:  os.Stdout,
		ContextDir:   contextDir,
	}
	return m.Client().BuildImage(opts)
}

func (m *Manager) Client() *docker.Client {
	return m.pool.Client
}

func (m *Manager) RunNode(node *Node) error {
	var resource *dockertest.Resource
	var err error

	if node.withRunOptions {
		resource, err = m.pool.RunWithOptions(node.RunOptions)
	} else {
		resource, err = m.pool.Run(node.repository, node.version, []string{})
	}

	if err != nil {
		stdOut, stdErr, _ := m.GetLogs(resource.Container.ID)
		return fmt.Errorf(
			"can't run container\n\n[error stream]:\n\n%s\n\n[output stream]:\n\n%s",
			stdErr,
			stdOut,
		)
	}

	// trying to get JSON-RPC server, to make sure node started properly
	// the last returned error before deadline exceeded will be returned from .Retry()
	err = m.pool.Retry(
		func() error {
			// recreating container instance because resource.Container.State
			// does not update properly by default
			c, err := m.Client().InspectContainer(resource.Container.ID)
			if err != nil {
				return fmt.Errorf("can't inspect container: %s", err.Error())
			}
			// if node failed to start, i.e. ExitCode != 0, return container logs
			if c.State.ExitCode != 0 {
				return fmt.Errorf(
					"can't start evmos node, container exit code: %d\n\n",
					c.State.ExitCode,
				)
			}
			// get host:port for current container in local network
			addr := resource.GetHostPort(jrpcPort + "/tcp")
			r, err := http.Get("http://" + addr)
			if err != nil {
				return fmt.Errorf("can't get node json-rpc server: %s", err)
			}
			defer r.Body.Close()
			return nil
		},
	)

	if err != nil {
		stdOut, stdErr, _ := m.GetLogs(resource.Container.ID)
		return fmt.Errorf(
			"can't start node: %s\n\n[error stream]:\n\n%s\n\n[output stream]:\n\n%s",
			err.Error(),
			stdErr,
			stdOut,
		)
	}
	m.CurrentNode = resource
	return nil
}

func (m *Manager) GetLogs(containerID string) (stdOut, stdErr string, err error) {
	var outBuf, errBuf bytes.Buffer
	opts := docker.LogsOptions{
		Container:    containerID,
		OutputStream: &outBuf,
		ErrorStream:  &errBuf,
		Stdout:       true,
		Stderr:       true,
	}
	err = m.Client().Logs(opts)
	if err != nil {
		return "", "", fmt.Errorf("can't get logs: %s", err)
	}
	return outBuf.String(), errBuf.String(), nil
}
