// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/rs/zerolog/log"
)

type TestCmd struct {
	Cover   bool `help:"Generate coverage report"`
	Race    bool `help:"Run race condition detector"`
	Profile bool `help:"Run with profiling"`
}

func (c *TestCmd) Run(ctx *Context) error {
	log.Print("Running test pipeline...")
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer client.Close()

	client = client.Pipeline("test", dagger.PipelineOpts{Description: "unit test pipeline"})

	hostSourceDir := client.Host().Directory("./", dagger.HostDirectoryOpts{
		Exclude: []string{"ci"},
	})
	source := client.Container().
		From(fmt.Sprintf("golang:%s-alpine", GO_VER)).
		WithMountedDirectory("/src", hostSourceDir)

	runner := source.WithWorkdir("/src")

	resultPath := "./results"

	testCmd := []string{"go", "test", "./...", "-shuffle=on"}

	if c.Race {
		runner = runner.
			WithEnvVariable("CGO_ENABLED", "1").
			WithExec([]string{"apk", "add", "build-base"})
		testCmd = append(testCmd, "-race")
	}

	if c.Cover {
		testCmd = append(testCmd, "-cover", "-covermode=atomic", "-coverprofile", filepath.Join(resultPath, "cover.out"))
	}

	if c.Profile {
		testCmd = append(testCmd, "-cpuprofile", filepath.Join(resultPath, "cpu.out"), "-memprofile", filepath.Join(resultPath, "mem.out"))
	}

	runner = runner.WithExec([]string{"mkdir", "-p", resultPath}).WithExec(testCmd)

	resultDir := runner.Directory(resultPath)

	_, err = resultDir.Export(ctx, resultPath)
	if err != nil {
		log.Err(err).Send()
		return err
	}

	entries, err := resultDir.Entries(ctx)
	if err != nil {
		log.Err(err).Send()
		return err
	}

	log.Printf("result dir contents: %s", entries)

	return nil
}
