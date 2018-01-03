package awslogin

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	_, err := NewSession("default")
	var expected error
	if err != expected {
		t.Errorf("expected %v to eq %v", err, expected)
	}
}

func TestNewCredentials(t *testing.T) {
	// write test
}

func TestBuildAssumeRoleProvider(t *testing.T) {
	// write test
}

func TestBuildSigninTokenRequestURL(t *testing.T) {
	u := BuildSigninTokenRequestURL("faderated-token")
	var expected string = "https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=faderated-token&SessionType=json"
	if u != expected {
		t.Errorf("expected %v to eq %v", u, expected)
	}
}

func TestRequestSigninToken(t *testing.T) {
	// write test
}

func TestBuildFederatedSession(t *testing.T) {
	// write test
}

func TestBuildSigninURL(t *testing.T) {
	u := BuildSigninURL("siginin-token")
	var expected string = "https://signin.aws.amazon.com/federation?Action=login&Destination=https%3A%2F%2Fconsole.aws.amazon.com%2F&Issuer=https%3A%2F%2Fgithub.com%2Fyouyo%2Fawslogin%2F&SigninToken=siginin-token"
	if u != expected {
		t.Errorf("expected %v to eq %v", u, expected)
	}
}
