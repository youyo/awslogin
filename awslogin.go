package awslogin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	// SigninBaseURL is request endpoint
	SigninBaseURL string = "https://signin.aws.amazon.com/federation"
)

type (
	TemporaryCredentials struct {
		SessionID    string `json:"sessionId"`
		SessionKey   string `json:"sessionKey"`
		SessionToken string `json:"sessionToken"`
	}

	SigninToken struct {
		Token string `json:"SigninToken"`
	}
)

// NewAwsSession
func NewAwsSession(profile string, durationSeconds time.Duration) (sess *session.Session) {
	sess = session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState:       session.SharedConfigEnable,
				Profile:                 profile,
				AssumeRoleDuration:      durationSeconds,
				AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			},
		),
	)
	return sess
}

func BuildTemporaryCredentials(sess *session.Session) (temporaryCredentials string, err error) {
	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		return "", err
	}

	tempCreds := &TemporaryCredentials{
		SessionID:    creds.AccessKeyID,
		SessionKey:   creds.SecretAccessKey,
		SessionToken: creds.SessionToken,
	}

	tempCredsByte, err := json.Marshal(*tempCreds)
	if err != nil {
		return "", err
	}

	return string(tempCredsByte), nil
}

func BuildSigninTokenRequestURL(temporaryCredentials, durationSeconds string) (requestUrl string) {
	values := url.Values{}
	values.Add("Action", "getSigninToken")
	values.Add("SessionType", "json")
	values.Add("Session", temporaryCredentials)
	values.Add("SessionDuration", durationSeconds)
	requestUrl = SigninBaseURL + "?" + values.Encode()
	return requestUrl
}

func RequestSigninToken(requestUrl string) (signinToken string, err error) {
	resp, err := http.Get(requestUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var st SigninToken
	if err = json.Unmarshal(body, &st); err != nil {
		return "", err
	}

	signinToken = st.Token

	return signinToken, nil
}

func BuildSigninURL(signinToken, region string) (signinUrl string) {
	values := url.Values{}
	values.Add("Action", "login")
	values.Add("Issuer", "https://github.com/youyo/awslogin/")
	values.Add("Destination", "https://"+region+".console.aws.amazon.com/")
	values.Add("SigninToken", signinToken)
	signinUrl = SigninBaseURL + "?" + values.Encode()

	return signinUrl
}
