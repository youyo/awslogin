package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/Songmu/prompter"
	"github.com/youyo/awslogin/lib/awslogin"
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
		profileName string
		list        bool
		version     bool
		output      bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.StringVar(&profileName, "profile", "", "Use Profile name")
	flags.StringVar(&profileName, "p", "", "Use Profile name. (Short)")
	flags.BoolVar(&list, "list", false, "Print available ARN list and quit.")
	flags.BoolVar(&list, "l", false, "Print available ARN list and quit. (Short)")
	flags.BoolVar(&output, "output", false, "Print login url and quiet. It will not login automatically.")
	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	flags.BoolVar(&version, "v", false, "Print version information and quit. (Short)")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		outdate, currentVersion, err := awslogin.VersionCheck(Version)
		if err != nil {
			fmt.Fprintf(cli.errStream, "%v\n", err)
		} else {
			if outdate {
				fmt.Fprintf(cli.errStream, "%s is not latest, you should upgrade to %s\n", Version, currentVersion)
			}
		}
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	//Load config
	cfg, err := awslogin.NewConfig()
	if err != nil {
		fmt.Fprintf(cli.errStream, "Could not load config. %s\n", err)
		return ExitCodeError
	}

	// Show Available ARNs
	if list {
		arns := cfg.AvailableArn()
		for _, v := range arns {
			fmt.Fprintf(cli.outStream, "%s\n", v)
		}
		return ExitCodeOK
	}

	// Check args
	if !awslogin.CheckArgProfileName(profileName) {
		fmt.Fprintf(cli.errStream, "The argument of profile name must not be empty.\n")
		return ExitCodeError
	}

	// Set profile name
	cfg.SetProfileName(profileName)

	// eFtch ARNs
	err = cfg.FetchArn()
	if err != nil {
		fmt.Fprintf(cli.errStream, "Could not fetch Arn.\n")
		return ExitCodeError
	}

	// If use to MFA, please enter mfa code.
	mfaCode := func() string {
		if cfg.MfaSerial != "" {
			return prompter.Prompt("Enter MFA code: ", "")
		}
		return ""
	}()

	// Set mfacode
	cfg.SetMfaCode(mfaCode)

	// Authentication by selected source profile
	s, err := awslogin.NewSession(cfg.SourceProfile)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Could not start new session.\n")
		return ExitCodeError
	}
	svc := awslogin.NewService(s)

	// Assume role
	resp, err := svc.AssumingRole(cfg)
	if err != nil {
		fmt.Fprintf(cli.errStream, "%s.\n", err)
		return ExitCodeError
	}

	// Build Federated Session
	fs, err := awslogin.BuildFederatedSession(resp)
	if err != nil {
		fmt.Fprintf(cli.errStream, "%s.\n", err)
		return ExitCodeError
	}

	// Request Signin Token
	url := awslogin.BuildSigninTokenRequestURL(fs)
	st, err := awslogin.RequestSigninToken(url)
	if err != nil {
		fmt.Fprintf(cli.errStream, "%s.\n", err)
		return ExitCodeError
	}

	// Build signin url
	url = awslogin.BuildSigninURL(st)

	// Open browser or output
	if output {
		fmt.Fprintf(cli.outStream, "%s\n", url)
		return ExitCodeOK
	} else {
		err = awslogin.Browsing(url)
		if err != nil {
			fmt.Fprintf(cli.errStream, "%s.\n", err)
			return ExitCodeError
		}
	}

	return ExitCodeOK
}
