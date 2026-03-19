package main

import (
	"github.com/alecthomas/kong"
	"github.com/youyo/awslogin/cmd"
)

// goreleaser ldflags で埋め込まれる
var version = "dev"

func main() {
	var cli cmd.CLI
	cli.AppVersion = version

	ctx := kong.Parse(&cli,
		kong.Name("awslogin"),
		kong.Description("Generate AWS Management Console login URL."),
		kong.UsageOnError(),
	)
	ctx.FatalIfErrorf(ctx.Run(&cli.Globals))
}
