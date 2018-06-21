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

type Signin struct {
	Token string `json:"SigninToken"`
}

type Awslogin struct {
	SourceProfile         string                             `json:"-"`
	ARN                   string                             `json:"-"`
	RoleSessionName       string                             `json:"-"`
	MfaSerial             string                             `json:"-"`
	DurationSeconds       int                                `json:"-"`
	AssumeRoleProvider    func(*stscreds.AssumeRoleProvider) `json:"-"`
	FederatedSessionID    string                             `json:"sessionId"`
	FederatedSessionKey   string                             `json:"sessionKey"`
	FederatedSessionToken string                             `json:"sessionToken"`
	FederatedSession      string                             `json:"-"`
	SigninTokenRequestURL string                             `json:"-"`
	SigninToken           string                             `json:"-"`
	SigninUrl             string                             `json:"-"`
}

func NewAwslogin(c *Config) *Awslogin {
	return &Awslogin{
		SourceProfile:   c.SourceProfile,
		ARN:             c.ARN,
		RoleSessionName: c.RoleSessionName,
		MfaSerial:       c.MfaSerial,
		DurationSeconds: c.DurationSeconds,
	}
}

func (al *Awslogin) BuildAssumeRoleProvider() {
	al.AssumeRoleProvider = func(p *stscreds.AssumeRoleProvider) {
		p.Duration = time.Duration(al.DurationSeconds) * time.Second
		p.RoleSessionName = al.RoleSessionName
		if al.MfaSerial != "" {
			p.SerialNumber = aws.String(al.MfaSerial)
			p.TokenProvider = stscreds.StdinTokenProvider
		}
	}
}

func (al *Awslogin) GetCredentials() error {
	cred := credentials.NewSharedCredentials("", al.SourceProfile)
	s, err := session.NewSession(&aws.Config{Credentials: cred})
	if err != nil {
		return err
	}
	creds, err := stscreds.NewCredentials(s, al.ARN, al.AssumeRoleProvider).Get()
	if err != nil {
		return err
	}
	al.FederatedSessionID = creds.AccessKeyID
	al.FederatedSessionKey = creds.SecretAccessKey
	al.FederatedSessionToken = creds.SessionToken
	return nil
}

func (al *Awslogin) GetFederatedSession() error {
	bytes, err := json.Marshal(al)
	if err != nil {
		return err
	}
	al.FederatedSession = string(bytes)
	return nil
}

func (al *Awslogin) BuildSigninTokenRequestURL() {
	values := url.Values{}
	values.Add("Action", "getSigninToken")
	values.Add("SessionType", "json")
	values.Add("Session", al.FederatedSession)
	al.SigninTokenRequestURL = SigninBaseURL + "?" + values.Encode()
}

func (al *Awslogin) RequestSigninToken() error {
	resp, err := http.Get(al.SigninTokenRequestURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var signin Signin
	if err = json.Unmarshal(body, &signin); err != nil {
		return err
	}
	al.SigninToken = signin.Token
	return nil
}

func (al *Awslogin) BuildSigninURL() {
	values := url.Values{}
	values.Add("Action", "login")
	values.Add("Issuer", "https://github.com/youyo/awslogin/")
	values.Add("Destination", "https://console.aws.amazon.com/")
	values.Add("SigninToken", al.SigninToken)
	al.SigninUrl = SigninBaseURL + "?" + values.Encode()
}

func (al *Awslogin) GetSigninUrl() string {
	return al.SigninUrl
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
