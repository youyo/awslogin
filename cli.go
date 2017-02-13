package main

import (
	"flag"
	"fmt"
	"io"

	latest "github.com/tcnksm/go-latest"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		rolename string
		config   string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&rolename, "rolename", "sample_role", "Use IAM-Role.")
	flags.StringVar(&rolename, "r", "sample_role", "Use IAM-Role. (Short)")

	flags.StringVar(&config, "config", "~/.awslogin/config.toml", "Use config file.")
	flags.StringVar(&config, "c", "~/.awslogin/config.toml", "Use config file. (Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	flags.BoolVar(&version, "v", false, "Print version information and quit. (Short)")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		versionCheck()
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	_ = rolename

	_ = config

	return ExitCodeOK
}

func versionCheck() {
	githubTag := &latest.GithubTag{
		Owner:      "youyo",
		Repository: "awslogin",
	}
	res, err := latest.Check(githubTag, Version)
	if err == nil {
		if res.Outdated {
			fmt.Printf("%s is not latest, you should upgrade to %s\n", Version, res.Current)
		}
	} else {
		fmt.Printf("Network is not unreachable. Can not check version.\n")
	}
}
