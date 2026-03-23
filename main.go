package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/youyo/awslogin/cmd"
	"github.com/youyo/awslogin/internal/jsonout"
)

// goreleaser ldflags で埋め込まれる
var version = "dev"

func main() {
	out := jsonout.New(os.Stdout, os.Stderr)

	var cli cmd.CLI
	cli.AppVersion = version

	parser, err := kong.New(&cli,
		kong.Name("awslogin"),
		kong.Description("Generate AWS Management Console login URL."),
		kong.UsageOnError(),
	)
	if err != nil {
		_ = out.WriteError(&jsonout.AppError{Code: jsonout.ErrInternal, Message: err.Error()})
		os.Exit(1)
	}

	kctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		_ = out.WriteError(&jsonout.AppError{Code: jsonout.ErrInvalidArgs, Message: err.Error()})
		os.Exit(1)
	}

	kctx.Bind(out)

	if err := kctx.Run(&cli.Globals, out); err != nil {
		_ = out.WriteError(err)
		os.Exit(1)
	}
}
