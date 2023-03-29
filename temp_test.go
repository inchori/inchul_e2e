package inchori_e2e

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"testing"

	"github.com/cosmos/go-bip39"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func TestDocker(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	buildOpts := docker.BuildImageOptions{
		Dockerfile:   "./DockerFile",
		ContextDir:   ".",
		NoCache:      true,
		OutputStream: io.Discard,
		ErrorStream:  os.Stdout,
		Name:         fmt.Sprintf("%s:%s", "dep_core", "0.1"),
	}
	err = pool.Client.BuildImage(buildOpts)
	require.NoError(t, err)

	var buyerMnemonic string
	bEntropy, err := bip39.NewEntropy(256)
	buyerMnemonic, err = bip39.NewMnemonic(bEntropy)
	require.NoError(t, err)

	var oracleMnemonic string
	oEntropy, err := bip39.NewEntropy(256)
	oracleMnemonic, err = bip39.NewMnemonic(oEntropy)
	require.NoError(t, err)

	runOpts := &dockertest.RunOptions{
		Repository: "dep_core",
		Tag:        "0.1",
		Name:       "panacea_dep_e2e",
		Env: []string{
			fmt.Sprintf("E2E_DATA_BUYER_MNEMONIC=%s", buyerMnemonic),
			fmt.Sprintf("E2E_ORACLE_MNEMONIC=%s", oracleMnemonic),
		},
	}

	panacea, err := pool.RunWithOptions(runOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})

	cmd := []string{"echo", "$E2E_DATA_BUYER_MNEMONIC"}
	var stdout4 bytes.Buffer
	exitCode, err := panacea.Exec(cmd, dockertest.ExecOptions{StdOut: &stdout4})
	require.NoError(t, err)
	require.Zero(t, exitCode)
	fmt.Println(stdout4.String())

}

func getHostPort(resource *dockertest.Resource, id string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		return resource.GetHostPort(id)
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}
	return u.Hostname()
}
