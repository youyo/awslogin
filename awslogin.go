package awslogin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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
func NewAwsSession(profile string, cache bool, cachePath string) (sess *session.Session, err error) {
	if cache {
		sess, err = newAwsSessionWithCache(profile, cachePath)
		if err != nil {
			return nil, err
		}

	} else {
		sess, err = newAwsSession(profile)
		if err != nil {
			return nil, err
		}
	}

	return sess, nil
}

func newAwsSession(profile string) (sess *session.Session, err error) {
	sess, err = session.NewSessionWithOptions(
		session.Options{
			SharedConfigState:       session.SharedConfigEnable,
			Profile:                 profile,
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		},
	)
	return sess, err
}

func newAwsSessionWithCreds(profile string, creds *credentials.Credentials) (sess *session.Session) {
	sess = session.Must(
		session.NewSessionWithOptions(
			session.Options{
				Config: aws.Config{
					Credentials: creds,
				},
				SharedConfigState:       session.SharedConfigEnable,
				Profile:                 profile,
				AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			},
		),
	)
	return sess
}

func newAwsSessionWithCache(profile, cachePath string) (sess *session.Session, err error) {
	c, err := NewCache(cachePath, profile)
	if err != nil {
		return nil, err
	}

	if cachedCreds, err := c.Load(); err != nil {
		sess, err = newAwsSession(profile)
		if err != nil {
			return nil, err
		}

		creds, err := sess.Config.Credentials.Get()
		if err != nil {
			return nil, err
		}

		c.Save(&creds)
	} else {
		creds := credentials.NewStaticCredentialsFromCreds(*cachedCreds)
		sess = newAwsSessionWithCreds(profile, creds)
	}

	return sess, nil
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

func BuildSigninTokenRequestURL(temporaryCredentials string) (requestUrl string) {
	values := url.Values{}
	values.Add("Action", "getSigninToken")
	values.Add("SessionType", "json")
	values.Add("Session", temporaryCredentials)
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

func BuildSigninURL(signinToken string) (signinUrl string) {
	values := url.Values{}
	values.Add("Action", "login")
	values.Add("Issuer", "https://github.com/youyo/awslogin/")
	values.Add("Destination", "https://console.aws.amazon.com/")
	values.Add("SigninToken", signinToken)
	values.Add("SessionDuration", "43200")
	signinUrl = SigninBaseURL + "?" + values.Encode()

	return signinUrl
}
