package awslogin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	SigninBaseURL string = "https://signin.aws.amazon.com/federation"
)

type (
	federatedSession struct {
		SessionID    string `json:"sessionId"`
		SessionKey   string `json:"sessionKey"`
		SessionToken string `json:"sessionToken"`
	}

	signinToken struct {
		Token string `json:"SigninToken"`
	}
)

func NewSession(sourceProfile string) (s *session.Session, err error) {
	cred := credentials.NewSharedCredentials("", sourceProfile)
	s, err = session.NewSession(&aws.Config{Credentials: cred})
	return
}

func NewCredentials(sess *session.Session, arn, roleSessionName, mfaSerial string, durationSeconds int) (creds credentials.Value, err error) {
	assumeRoleProvider := buildAssumeRoleProvider(roleSessionName, mfaSerial, durationSeconds)
	creds, err = stscreds.NewCredentials(sess, arn, assumeRoleProvider).Get()
	return
}

func buildAssumeRoleProvider(roleSessionName, mfaSerial string, durationSeconds int) (f func(p *stscreds.AssumeRoleProvider)) {
	f = func(p *stscreds.AssumeRoleProvider) {
		p.Duration = time.Duration(durationSeconds) * time.Second
		p.RoleSessionName = roleSessionName
		if mfaSerial != "" {
			p.SerialNumber = aws.String(mfaSerial)
			p.TokenProvider = stscreds.StdinTokenProvider
		}
	}
	return
}

func BuildSigninTokenRequestURL(fs string) (u string) {
	values := url.Values{}
	values.Add("Action", "getSigninToken")
	values.Add("SessionType", "json")
	values.Add("Session", fs)
	u = SigninBaseURL + "?" + values.Encode()
	return
}

func RequestSigninToken(url string) (st string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var ST signinToken
	if err = json.Unmarshal(body, &ST); err != nil {
		return
	}
	st = ST.Token
	return
}

func BuildFederatedSession(accessKeyId, secretAccessKey, sessionToken string) (j string, err error) {
	fs := &federatedSession{
		SessionID:    accessKeyId,
		SessionKey:   secretAccessKey,
		SessionToken: sessionToken,
	}
	b, err := json.Marshal(*fs)
	j = string(b)
	return
}

func BuildSigninURL(st string) (u string) {
	values := url.Values{}
	values.Add("Action", "login")
	values.Add("Issuer", "https://github.com/youyo/awslogin/")
	values.Add("Destination", "https://console.aws.amazon.com/")
	values.Add("SigninToken", st)
	u = SigninBaseURL + "?" + values.Encode()
	return
}
