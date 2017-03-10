package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"github.com/skratchdot/open-golang/open"
	latest "github.com/tcnksm/go-latest"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota

	SigninBaseURL string = "https://signin.aws.amazon.com/federation"
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

type (
	service struct {
		*sts.STS
	}

	federatedSession struct {
		SessionID    string `json:"sessionId"`
		SessionKey   string `json:"sessionKey"`
		SessionToken string `json:"sessionToken"`
	}

	signinToken struct {
		Token string `json:"SigninToken"`
	}
)

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		rolename string
		list     bool
		version  bool
		profile  string
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.StringVar(&rolename, "rolename", "", "Use IAM-Role.")
	flags.StringVar(&rolename, "r", "", "Use IAM-Role. (Short)")
	flags.BoolVar(&list, "list", false, "Print available ARN list and quit.")
	flags.BoolVar(&list, "l", false, "Print available ARN list and quit. (Short)")
	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	flags.BoolVar(&version, "v", false, "Print version information and quit. (Short)")
	flags.StringVar(&profile, "profile", "default", "Use default profile.")

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

	//Load config
	cfgPath := configPath()
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Could not load config. %s\n", err)
		return ExitCodeError
	}

	// Show Available ARNs
	if list {
		l := availableArn(cfg)
		for _, v := range l {
			fmt.Fprintf(cli.outStream, "%s\n", v)
		}
		return ExitCodeOK
	}

	if !checkArgRoleName(rolename) {
		fmt.Fprintf(cli.errStream, "The argument of rolename must not be empty.\n")
		return ExitCodeError
	}

	// new service
	s, err := newSession(profile)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Could not start new session.\n")
		return ExitCodeError
	}
	svc := newService(s)

	// fetch arn
	arn, err := fetchArn(cfg, rolename)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Could not fetch Arn.\n")
		return ExitCodeError
	}

	// assume role
	resp, err := svc.assumeRole(rolename, arn)
	if err != nil {
		fmt.Fprintf(cli.errStream, "%s.\n", err)
		return ExitCodeError
	}

	// build Federated Session
	fs, err := buildFederatedSession(resp)

	// request Signin Token
	url := buildSigninTokenRequestURL(fs)
	st, err := requestSigninToken(url)
	if err != nil {
		fmt.Fprintf(cli.errStream, "%s.\n", err)
		return ExitCodeError
	}

	// build signin url
	url = buildSigninURL(st)

	// open browse
	err = open.Start(url)
	if err != nil {
		fmt.Fprintf(cli.errStream, "%s.\n", err)
		return ExitCodeError
	}

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

func newSession(p string) (s *session.Session, err error) {
	cred := credentials.NewSharedCredentials("", p)
	s, err = session.NewSession(&aws.Config{Credentials: cred})
	return
}

func newService(s *session.Session) (svc *service) {
	svc = &service{sts.New(s)}
	return
}

func (svc *service) assumeRole(roleName, arn string) (resp *sts.AssumeRoleOutput, err error) {
	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(roleName),
	}
	return svc.AssumeRole(params)
}

func buildSigninTokenRequestURL(fs string) (u string) {
	values := url.Values{}
	values.Add("Action", "getSigninToken")
	values.Add("SessionType", "json")
	values.Add("Session", fs)
	u = SigninBaseURL + "?" + values.Encode()
	return
}

func requestSigninToken(url string) (st string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ST signinToken
	json.Unmarshal(body, &ST)
	return ST.Token, nil
}

func configPath() (c string) {
	c = filepath.Join(os.Getenv("HOME"), ".aws/config")
	return
}

func loadConfig(configPath string) (cfg *ini.File, err error) {
	cfg, err = ini.Load(configPath)
	return
}

func availableArn(cfg *ini.File) (list []string) {
	sections := cfg.Sections()
	for _, s := range sections {
		if s.HasKey("role_arn") {
			n := strings.Replace(s.Name(), "profile ", "", 1)
			list = append(list, n)
		}
	}
	return
}

func fetchArn(cfg *ini.File, roleName string) (arn string, err error) {
	s := "profile " + roleName
	arn = cfg.Section(s).Key("role_arn").String()
	if arn == "" {
		return "", errors.New("Could not fetch Arn")
	}
	return
}

func buildFederatedSession(resp *sts.AssumeRoleOutput) (j string, err error) {
	fs := &federatedSession{
		SessionID:    *resp.Credentials.AccessKeyId,
		SessionKey:   *resp.Credentials.SecretAccessKey,
		SessionToken: *resp.Credentials.SessionToken,
	}
	b, err := json.Marshal(*fs)
	j = string(b)
	return
}

func buildSigninURL(st string) (u string) {
	values := url.Values{}
	values.Add("Action", "login")
	values.Add("Issuer", "https://github.com/youyo/awslogin/")
	values.Add("Destination", "https://console.aws.amazon.com/")
	values.Add("SigninToken", st)
	u = SigninBaseURL + "?" + values.Encode()
	return
}

func checkArgRoleName(r string) bool {
	if r == "" {
		return false
	}
	return true
}
