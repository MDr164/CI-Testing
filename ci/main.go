// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PROJ_NAME      = "template"
	GO_VER         = "1.20"
	GORELEASER_VER = "1.16.1"
)

type Context struct {
	context.Context
}

var cli struct {
	Test    TestCmd    `cmd:"" help:"Run test pipeline"`
	Build   BuildCmd   `cmd:"" help:"Run build pipeline"`
	Release ReleaseCmd `cmd:"" help:"Run release pipeline"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	ctx := kong.Parse(&cli)
	if err := ctx.Run(&Context{context.Background()}); err != nil {
		log.Fatal().Err(err).Send()
	}
}
