package inchori_e2e

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/go-bip39"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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
		Hostname: "panacea",
	}

	panacea, err := pool.RunWithOptions(runOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})

	time.Sleep(30 * time.Second)

	cmd := []string{"panacead", "q", "block", "1"}
	var stdout2 bytes.Buffer
	exitCode1, err := panacea.Exec(cmd, dockertest.ExecOptions{StdOut: &stdout2})
	require.NoError(t, err)
	require.Zero(t, exitCode1)

	var resultBlock ctypes.ResultBlock
	err = legacy.Cdc.UnmarshalJSON(stdout2.Bytes(), &resultBlock)
	require.NoError(t, err)
	fmt.Println(resultBlock.BlockID.Hash.String())

	oracleBuildOpts := docker.BuildImageOptions{
		Dockerfile:   "./oracle.DockerFile",
		ContextDir:   ".",
		NoCache:      true,
		OutputStream: io.Discard,
		ErrorStream:  os.Stdout,
		Name:         fmt.Sprintf("%s:%s", "oracle_init", "0.1"),
	}
	err = pool.Client.BuildImage(oracleBuildOpts)
	require.NoError(t, err)

	oracleRunOpts := &dockertest.RunOptions{
		Repository: "oracle_init",
		Tag:        "0.1",
		Name:       "oracle_dep_e2e",
		Env: []string{
			fmt.Sprintf("CHAIN_ID=%s", "testing"),
			fmt.Sprintf("ORACLE_MNEMONIC=%s", oracleMnemonic),
			fmt.Sprintf("TRUSTED_BLOCK_HASH=%s", resultBlock.BlockID.Hash.String()),
		},
	}

	oracle, err := pool.RunWithOptions(oracleRunOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})

	cmd4 := []string{"cat", "/oracle/.oracle/oracle_pub_key.json"}
	var stdout bytes.Buffer
	exitCode, err := oracle.Exec(cmd4, dockertest.ExecOptions{StdOut: &stdout})
	require.Zero(t, exitCode)
	require.NoError(t, err)
	fmt.Println(stdout.String())
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

	return u.Hostname() + ":" + resource.GetPort(id)
}
