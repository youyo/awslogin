package signin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	// SigninBaseURL は AWS Federation API のエンドポイント
	SigninBaseURL = "https://signin.aws.amazon.com/federation"

	// DefaultRegion はリージョン未設定時のフォールバック値
	DefaultRegion = "us-east-1"
)

// TemporaryCredentials は Federation API に渡す一時認証情報の JSON 構造体
type TemporaryCredentials struct {
	SessionID    string `json:"sessionId"`
	SessionKey   string `json:"sessionKey"`
	SessionToken string `json:"sessionToken"`
}

// signinTokenResponse は Federation API の getSigninToken レスポンス
type signinTokenResponse struct {
	Token string `json:"SigninToken"`
}

// BuildTemporaryCredentials は認証情報を Federation API 用の JSON 文字列に変換する
func BuildTemporaryCredentials(accessKeyID, secretAccessKey, sessionToken string) (string, error) {
	creds := &TemporaryCredentials{
		SessionID:    accessKeyID,
		SessionKey:   secretAccessKey,
		SessionToken: sessionToken,
	}

	data, err := json.Marshal(creds)
	if err != nil {
		return "", fmt.Errorf("failed to marshal credentials: %w", err)
	}

	return string(data), nil
}

// BuildSigninTokenRequestURL は SigninToken 取得用の URL を構築する
func BuildSigninTokenRequestURL(credentials string, durationSeconds int) string {
	values := url.Values{}
	values.Add("Action", "getSigninToken")
	values.Add("SessionType", "json")
	values.Add("Session", credentials)
	values.Add("SessionDuration", strconv.Itoa(durationSeconds))

	return SigninBaseURL + "?" + values.Encode()
}

// RequestSigninToken は Federation API から SigninToken を取得する
// HTTP ステータスコード 200 以外の場合はエラーを返す
func RequestSigninToken(requestURL string) (string, error) {
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", fmt.Errorf("failed to request signin token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("federation API returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var st signinTokenResponse
	if err := json.Unmarshal(body, &st); err != nil {
		return "", fmt.Errorf("failed to parse signin token response: %w", err)
	}

	return st.Token, nil
}

// consoleURL はリージョンに応じたコンソール URL を返す
func consoleURL(region string) string {
	if region == DefaultRegion {
		return "https://console.aws.amazon.com/"
	}
	return "https://" + region + ".console.aws.amazon.com/"
}

// BuildSigninURL はログイン URL を構築する
// 内部で ResolveRegion を呼び、空文字リージョンも自動的に us-east-1 にフォールバック
func BuildSigninURL(signinToken, region string) string {
	resolved := ResolveRegion(region)

	values := url.Values{}
	values.Add("Action", "login")
	values.Add("Issuer", "https://github.com/youyo/awslogin/")
	values.Add("Destination", consoleURL(resolved))
	values.Add("SigninToken", signinToken)

	return SigninBaseURL + "?" + values.Encode()
}

// ResolveRegion はリージョンを解決する。空文字の場合 us-east-1 にフォールバックする。
func ResolveRegion(region string) string {
	if region == "" {
		return DefaultRegion
	}
	return region
}
