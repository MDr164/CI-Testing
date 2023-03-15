// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/rs/zerolog/log"
)

type ReleaseCmd struct{}

func (c *ReleaseCmd) Run(ctx *Context) error {
	log.Print("Running release pipeline...")
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer client.Close()

	client = client.Pipeline("release", dagger.PipelineOpts{Description: "release pipeline using GoReleaser"})

	hostSourceDir := client.Host().Directory("./", dagger.HostDirectoryOpts{
		Exclude: []string{"ci"},
	})
	source := client.Container().
		From(fmt.Sprintf("golang:%s-alpine", GO_VER)).
		WithMountedDirectory("/src", hostSourceDir)

	runner := source.
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "go", "cosign", "syft", "gpg", "curl", "bash", "ca-certificates"}).
		WithExec([]string{"curl", "-fL", fmt.Sprintf("https://github.com/goreleaser/goreleaser/releases/download/v%s/goreleaser_%s_x86_64.apk", GORELEASER_VER, GORELEASER_VER), "-o", "/tmp/goreleaser.apk"}).
		WithExec([]string{"apk", "add", "--no-cache", "--allow-untrusted", "/tmp/goreleaser.apk"})

	resultPath := "./results"

	runner = runner.
		WithSecretVariable("GITHUB_TOKEN", client.Host().EnvVariable("GITHUB_TOKEN").Secret()).
		WithSecretVariable("COSIGN_KEY", client.Host().EnvVariable("COSIGN_KEY").Secret()).
		WithSecretVariable("COSIGN_PWD", client.Host().EnvVariable("COSIGN_PWD").Secret()).
		WithExec([]string{"goreleaser", "release", "--snapshot"})

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
