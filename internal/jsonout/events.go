package jsonout

// SSOSessionExpiredEvent は SSO セッション期限切れイベント
type SSOSessionExpiredEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// SSOAuthRequiredEvent はデバイス認証必要イベント
type SSOAuthRequiredEvent struct {
	Type                string `json:"type"`
	VerificationCode    string `json:"verification_code"`
	VerificationURL     string `json:"verification_url"`
	VerificationURLBase string `json:"verification_url_base,omitempty"`
}

// BrowserOpenedEvent はブラウザ起動成功イベント
type BrowserOpenedEvent struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// BrowserOpenFailedEvent はブラウザ起動失敗イベント
type BrowserOpenFailedEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// SSOAuthCompleteEvent は認証完了イベント
type SSOAuthCompleteEvent struct {
	Type string `json:"type"`
}

// NewSSOSessionExpiredEvent を作成
func NewSSOSessionExpiredEvent() SSOSessionExpiredEvent {
	return SSOSessionExpiredEvent{Type: "sso_session_expired", Message: "SSO session expired. Starting SSO login..."}
}

// NewSSOAuthRequiredEvent を作成
func NewSSOAuthRequiredEvent(code, url, urlBase string) SSOAuthRequiredEvent {
	return SSOAuthRequiredEvent{Type: "sso_auth_required", VerificationCode: code, VerificationURL: url, VerificationURLBase: urlBase}
}

// NewBrowserOpenedEvent を作成
func NewBrowserOpenedEvent(url string) BrowserOpenedEvent {
	return BrowserOpenedEvent{Type: "browser_opened", URL: url}
}

// NewBrowserOpenFailedEvent を作成
func NewBrowserOpenFailedEvent() BrowserOpenFailedEvent {
	return BrowserOpenFailedEvent{Type: "browser_open_failed", Message: "Could not open browser automatically. Please open the URL manually."}
}

// NewSSOAuthCompleteEvent を作成
func NewSSOAuthCompleteEvent() SSOAuthCompleteEvent {
	return SSOAuthCompleteEvent{Type: "sso_auth_complete"}
}
