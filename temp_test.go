package inchori_e2e

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func TestDocker(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	coreBuildArgs := map[string]string{"E2E_DATA_BUYER_MNEMONIC": "delay inherit claw novel phone fee truck fault isolate merry across north mosquito race judge rally unaware brush employ pencil half gossip ramp suspect", "E2E_ORACLE_MNEMONIC": "gasp shy describe man hello blossom motor monkey seven mule shallow almost bunker hello wife clarify tissue best actress hub wisdom crane ridge heavy"}
	buildArgs := make([]docker.BuildArg, 0, len(coreBuildArgs))
	for k, v := range coreBuildArgs {
		bArg := docker.BuildArg{
			Name:  k,
			Value: v,
		}
		buildArgs = append(buildArgs, bArg)
	}

	buildOpts := docker.BuildImageOptions{
		Dockerfile:   "./DockerFile",
		ContextDir:   ".",
		NoCache:      true,
		OutputStream: io.Discard,
		ErrorStream:  os.Stdout,
		Name:         fmt.Sprintf("%s:%s", "core_init", "0.1"),
		BuildArgs:    buildArgs,
	}
	err = pool.Client.BuildImage(buildOpts)
	require.NoError(t, err)

	runOpts := &dockertest.RunOptions{
		Repository: "core_init",
		Tag:        "0.1",
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

	panaceaHost := getHostPort(panacea, "1317/tcp")

	require.NoError(t, err)
	defer panacea.Close()

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
		Env: []string{"CHAIN_ID=testing", "ORACLE_MNEMONIC=`gasp shy describe man hello blossom motor monkey seven mule shallow almost bunker hello wife clarify tissue best actress hub wisdom crane ridge heavy`",
			fmt.Sprintf("PANACEA_HOST=%s", panaceaHost)},
	}

	oracle, err := pool.RunWithOptions(oracleRunOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})

	cmd1 := []string{"echo", "$CHAIN_ID"}
	cmd2 := []string{"echo", "$ORACLE_MNEMONIC"}
	cmd3 := []string{"echo", "$PANACEA_HOST"}

	var stdout1 bytes.Buffer
	exitCode1, err := oracle.Exec(cmd1, dockertest.ExecOptions{StdOut: &stdout1})
	require.NoError(t, err)
	require.Zero(t, exitCode1)
	fmt.Println(stdout1.String())

	var stdout2 bytes.Buffer
	exitCode2, err := oracle.Exec(cmd2, dockertest.ExecOptions{StdOut: &stdout2})
	require.NoError(t, err)
	require.Zero(t, exitCode2)
	fmt.Println(stdout2.String())

	var stdout3 bytes.Buffer
	exitCode3, err := oracle.Exec(cmd3, dockertest.ExecOptions{StdOut: &stdout3})
	require.NoError(t, err)
	require.Zero(t, exitCode3)
	fmt.Println(stdout3.String())

	cmd4 := []string{"cat", "config.toml"}
	var stdout bytes.Buffer
	exitCode, err = oracle.Exec(cmd4, dockertest.ExecOptions{StdOut: &stdout})
	require.NoError(t, err)
	require.Zero(t, exitCode)
	fmt.Println(stdout.String())

	require.NoError(t, err)
	defer oracle.Close()
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
