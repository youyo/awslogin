package awslogin

import (
	"reflect"
	"testing"
)

func TestNewSession(t *testing.T) {
	_, err := NewSession("default")
	var expected error
	if err != expected {
		t.Errorf("expected %v to eq %v", err, expected)
	}
}

func TestNewService(t *testing.T) {
	s, _ := NewSession("default")
	svc := NewService(s)
	var expected *service
	if reflect.TypeOf(svc) != reflect.TypeOf(expected) {
		t.Errorf("expected %v to eq %v", reflect.TypeOf(svc), reflect.TypeOf(expected))
	}
}

func TestExtractRoleSessionName(t *testing.T) {
	roleSessionName := extractRoleSessionName("arn:aws:iam::xxxxxxxxxxxx:role/xxxxx")
	expected := "xxxxx"
	if roleSessionName != expected {
		t.Errorf("expected %v to eq %v", roleSessionName, expected)
	}
}

func TestBuildSigninTokenRequestURL(t *testing.T) {
	u := BuildSigninTokenRequestURL("faderated-token")
	expected := "https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=faderated-token&SessionType=json"
	if u != expected {
		t.Errorf("expected %v to eq %v", u, expected)
	}
}

func TestBuildSigninURL(t *testing.T) {
	u := BuildSigninURL("siginin-token")
	expected := "https://signin.aws.amazon.com/federation?Action=login&Destination=https%3A%2F%2Fconsole.aws.amazon.com%2F&Issuer=https%3A%2F%2Fgithub.com%2Fyouyo%2Fawslogin%2F&SigninToken=siginin-token"
	if u != expected {
		t.Errorf("expected %v to eq %v", u, expected)
	}
}
