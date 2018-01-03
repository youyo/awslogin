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
	Name       string
	Version    string
	CommitHash string
	BuildTime  string
	GoVersion  string
	app        string
)

const (
	Peco          string = "peco"
	PecoGithubUrl string = "https://github.com/peco/peco"
)

var RootCmd = &cobra.Command{
	Use:   "awslogin",
	Short: "",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := execRoot(); err != nil {
			log.Fatal(err)
		}
	},
}

func execRoot() (err error) {
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

	profile, err := cmd.Output()
	if err != nil {
		errorMessage := fmt.Sprintf("'%s' is Required command. Please install it. %s ", Peco, PecoGithubUrl)
		err = errors.Wrap(err, errorMessage)
		return
	}

	cfg, err := awslogin.NewConfig()
	if err != nil {
		return
	}

	cfg.SetProfileName(strings.TrimRight(string(profile), "\n"))

	if err = cfg.FetchArn(); err != nil {
		return
	}

	sess, err := awslogin.NewSession(cfg.SourceProfile)
	if err != nil {
		return
	}

	creds, err := awslogin.NewCredentials(sess, cfg.ARN, cfg.RoleSessionName, cfg.MfaSerial)
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

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().StringVarP(&app, "app", "a", "", "Opens with the specified application.")
}

func initConfig() {}
