package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
func TestRun_versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./awslogin -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("awslogin version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}
*/

func TestCheckArgRoleName(t *testing.T) {
	b := checkArgRoleName("")
	expected := false
	if b {
		t.Errorf("expected %v to eq %v", b, expected)
	}
}

func TestBuildSigninURL(t *testing.T) {
	u := buildSigninURL("siginin-token")
	expected := "https://signin.aws.amazon.com/federation?Action=login&Destination=https%3A%2F%2Fconsole.aws.amazon.com%2F&Issuer=https%3A%2F%2Fgithub.com%2Fyouyo%2Fawslogin%2F&SigninToken=siginin-token"
	if u != expected {
		t.Errorf("expected %v to eq %v", u, expected)
	}
}

func TestNewSession(t *testing.T) {
	_, err := newSession("default")
	var expected error
	if err != nil {
		t.Errorf("expected %v to eq %v", err, expected)
	}
}

func TestBuildSigninTokenRequestURL(t *testing.T) {
	u := buildSigninTokenRequestURL("faderated-token")
	expected := "https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=faderated-token&SessionType=json"
	if u != expected {
		t.Errorf("expected %v to eq %v", u, expected)
	}
}

func TestConfigPath(t *testing.T) {
	c := configPath()
	expected := filepath.Join(os.Getenv("HOME"), ".aws/config")
	if c != expected {
		t.Errorf("expected %v to eq %v", c, expected)
	}
}

func TestLoadConfig(t *testing.T) {
	_, err := loadConfig("./tests/config")
	var expected error
	if err != nil {
		t.Errorf("expected %v to eq %v", err, expected)
	}
}

func TestAvailableArn(t *testing.T) {
	cfg, _ := loadConfig("./tests/config")
	list := availableArn(cfg)
	var expected []string
	expected = append(expected, "test")
	if list[0] != expected[0] {
		t.Errorf("expected %v to eq %v", list, expected)
	}
}

func TestFetchArn(t *testing.T) {
	cfg, _ := loadConfig("./tests/config")
	_, err := fetchArn(cfg, "test")
	var expected error
	if err != nil {
		t.Errorf("expected %v to eq %v", err, expected)
	}
}

func TestRun_listFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./awslogin -list", " ")
	expected := 0
	status := cli.Run(args)
	if status != expected {
		t.Errorf("expected %v to eq %v", status, expected)
	}
}
