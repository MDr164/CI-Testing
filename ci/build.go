// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/rs/zerolog/log"
)

type BuildCmd struct {
	Strip bool `help:"Strip resulting binary"`
}

func (c *BuildCmd) Run(ctx *Context) error {
	log.Print("Running build pipeline...")
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Err(err).Send()
		return err
	}
	defer client.Close()

	client = client.Pipeline("build", dagger.PipelineOpts{Description: "build pipeline for main source code"})

	hostSourceDir := client.Host().Directory("./", dagger.HostDirectoryOpts{
		Exclude: []string{"ci"},
	})
	source := client.Container().
		From(fmt.Sprintf("golang:%s-alpine", GO_VER)).
		WithMountedDirectory("/src", hostSourceDir)

	runner := source.WithWorkdir("/src")

	resultPath := "./build"

	buildCmd := []string{"go", "build"}

	if c.Strip {
		buildCmd = append(buildCmd, "-ldflags=-s -w")
	}

	platforms := []string{"linux", "darwin", "windows"}
	arches := []string{"amd64", "arm64"}

	for _, platform := range platforms {
		for _, arch := range arches {
			path := filepath.Join(resultPath, platform, arch, fmt.Sprintf("%s-%s-%s", PROJ_NAME, platform, arch))

			buildCmd = append(buildCmd, "-o", path)

			runner = runner.
				WithEnvVariable("GOOS", platform).
				WithEnvVariable("GOARCH", arch).
				WithExec(buildCmd)
		}
	}

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

	log.Printf("build dir contents: %s", entries)

	return nil
}
