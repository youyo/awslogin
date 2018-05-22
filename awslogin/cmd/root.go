package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/youyo/awslogin"
)

var (
	Name            string
	Version         string
	CommitHash      string
	BuildTime       string
	GoVersion       string
	app             string
	profile         string
	readFromEnv     bool
	durationSeconds int
)

const (
	Peco          string = "peco"
	PecoGithubUrl string = "https://github.com/peco/peco"
)

var RootCmd = &cobra.Command{
	Use:   "awslogin",
	Short: "Login to the AWS management console.",
	Long:  `Using AssumeRole, accept IAMRole and login to the AWS management console.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := execRoot(); err != nil {
			log.Fatal(err)
		}
	},
}

func execRoot() (err error) {
	p, err := loadProfile(profile)
	if err != nil {
		return
	}

	cfg, err := awslogin.NewConfig()
	if err != nil {
		return
	}

	cfg.SetProfileName(p)

	if err = cfg.FetchArn(); err != nil {
		return
	}

	if durationSeconds != 3600 {
		cfg.SetDurationSeconds(durationSeconds)
	}

	sess, err := awslogin.NewSession(cfg.SourceProfile)
	if err != nil {
		return
	}

	creds, err := awslogin.NewCredentials(sess, cfg.ARN, cfg.RoleSessionName, cfg.MfaSerial, cfg.DurationSeconds)
	if err != nil {
		return
	}

	federatedSession, err := awslogin.BuildFederatedSession(
		creds.AccessKeyID,
		creds.SecretAccessKey,
		creds.SessionToken,
	)
	if err != nil {
		return
	}

	signinTokenRequestUrl := awslogin.BuildSigninTokenRequestURL(federatedSession)
	signinToken, err := awslogin.RequestSigninToken(signinTokenRequestUrl)
	if err != nil {
		return
	}

	signinUrl := awslogin.BuildSigninURL(signinToken)

	// Open browser
	if app != "" {
		err = awslogin.BrowsingSpecificApp(signinUrl, app)
	} else {
		err = awslogin.Browsing(signinUrl)
	}
	return
}

func loadProfile(p string) (profile string, err error) {
	if p != "" {
		profile = p
		return
	}

	// Read from env-variable
	if readFromEnv {
		profile = os.Getenv("AWS_PROFILE")
		fmt.Fprintf(os.Stdout, "The selected profile is %s\n", profile)
		return
	}

	cmd := exec.Command(Peco)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	buf := &bytes.Buffer{}
	if err = execList(buf); err != nil {
		return
	}

	list := buf.String()
	io.WriteString(stdin, list)
	stdin.Close()

	byteProfile, err := cmd.Output()
	if err != nil {
		errorMessage := fmt.Sprintf("'%s' is Required command. Please install it. %s ", Peco, PecoGithubUrl)
		err = errors.Wrap(err, errorMessage)
		return
	}
	profile = strings.TrimRight(string(byteProfile), "\n")
	fmt.Fprintf(os.Stdout, "The selected profile is %s\n", profile)
	return
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().StringVarP(&app, "app", "a", "", "Opens with the specified application.")
	RootCmd.Flags().StringVarP(&profile, "profile", "p", "", "Use a specific profile.")
	RootCmd.Flags().BoolVarP(&readFromEnv, "read-from-env", "e", false, "Use a specific profile read from the environment. [$AWS_PROFILE]")
	RootCmd.Flags().IntVarP(&durationSeconds, "duration-seconds", "d", 3600, "Request a session duration seconds. 900 - 43200")
}

func initConfig() {}