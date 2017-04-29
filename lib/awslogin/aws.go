package awslogin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	SigninBaseURL string = "https://signin.aws.amazon.com/federation"
)

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

func NewSession(sourceProfile string) (s *session.Session, err error) {
	cred := credentials.NewSharedCredentials("", sourceProfile)
	s, err = session.NewSession(&aws.Config{Credentials: cred})
	return
}

func NewService(s *session.Session) (svc *service) {
	svc = &service{sts.New(s)}
	return
}

func (svc *service) AssumingRole(cfg *Config) (resp *sts.AssumeRoleOutput, err error) {
	roleSessionName := extractRoleSessionName(cfg.ARN)
	params := func() *sts.AssumeRoleInput {
		if cfg.MfaSerial == "" {
			return &sts.AssumeRoleInput{
				RoleArn:         aws.String(cfg.ARN),
				RoleSessionName: aws.String(roleSessionName),
			}
		}
		return &sts.AssumeRoleInput{
			RoleArn:         aws.String(cfg.ARN),
			RoleSessionName: aws.String(roleSessionName),
			SerialNumber:    aws.String(cfg.MfaSerial),
			TokenCode:       aws.String(cfg.MfaCode),
		}
	}()
	return svc.AssumeRole(params)
}

func extractRoleSessionName(arn string) (roleSessionName string) {
	roleSessionName = strings.Split(arn, "/")[1] + "@awslogin"
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

func BuildFederatedSession(resp *sts.AssumeRoleOutput) (j string, err error) {
	fs := &federatedSession{
		SessionID:    *resp.Credentials.AccessKeyId,
		SessionKey:   *resp.Credentials.SecretAccessKey,
		SessionToken: *resp.Credentials.SessionToken,
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
