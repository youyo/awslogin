package jsonout

// エラーコード定数
const (
	ErrInvalidArgs          = "INVALID_ARGS"
	ErrConfigLoadFailed     = "CONFIG_LOAD_FAILED"
	ErrSSOSessionExpired    = "SSO_SESSION_EXPIRED"
	ErrSSOLoginFailed       = "SSO_LOGIN_FAILED"
	ErrSSOLegacyConfig      = "SSO_LEGACY_CONFIG"
	ErrOIDCRegisterFailed   = "OIDC_REGISTER_FAILED"
	ErrOIDCDeviceAuthFailed = "OIDC_DEVICE_AUTH_FAILED"
	ErrOIDCTokenExpired     = "OIDC_TOKEN_EXPIRED"
	ErrFederationAPIFailed  = "FEDERATION_API_FAILED"
	ErrBrowserOpenFailed    = "BROWSER_OPEN_FAILED"
	ErrInternal             = "INTERNAL_ERROR"
)

// AppError は構造化エラーコードを持つアプリケーションエラー
type AppError struct {
	Code    string
	Message string
	Details string
	Cause   error
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Cause }
