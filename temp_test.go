package inchori_e2e

import (
	"fmt"
	"io"
	"os"
	"testing"

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
		Name:         fmt.Sprintf("%s:%s", "core_init", "0.1"),
	}
	err = pool.Client.BuildImage(buildOpts)
	require.NoError(t, err)

	runOpts := &dockertest.RunOptions{
		Repository: "core_init",
		Tag:        "0.1",
		Env: []string{"E2E_DATA_BUYER_MNEMONIC=`delay inherit claw novel phone fee truck fault isolate merry across north mosquito race judge rally unaware brush employ pencil half gossip ramp suspect`",
			"E2E_ORACLE_MNEMONIC=`gasp shy describe man hello blossom motor monkey seven mule shallow almost bunker hello wife clarify tissue best actress hub wisdom crane ridge heavy`"},
	}

	panacea, err := pool.RunWithOptions(runOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})

	require.NoError(t, err)
	defer panacea.Close()
}
