package signin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/youyo/awslogin/internal/jsonout"
	"github.com/youyo/awslogin/internal/sso"
)

// AWSCredentials は AWS 認証情報とリージョン情報を保持する
type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}

// loadCredentialsOnce は AWS SDK v2 で認証情報とリージョンを一度取得する内部関数。
func loadCredentialsOnce(ctx context.Context) (*AWSCredentials, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve AWS credentials: %w", err)
	}

	return &AWSCredentials{
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
		Region:          ResolveRegion(cfg.Region),
	}, nil
}

// LoadCredentials は AWS SDK v2 で認証情報とリージョンを取得する。
// AWS_PROFILE 環境変数、~/.aws/credentials、IAM ロール等、SDK の config loader が自動解決する。
// SSO セッションが期限切れ（InvalidGrantException）の場合は OIDC デバイス認証フローを起動して再試行する。
func LoadCredentials(ctx context.Context, events io.Writer) (*AWSCredentials, error) {
	result, err := loadCredentialsOnce(ctx)
	if err == nil {
		return result, nil
	}

	// InvalidGrantException でなければそのままエラーを返す
	if !sso.IsInvalidGrantError(err) {
		return nil, err
	}

	// SSO プロファイルか確認する
	ssoCfg, ssoErr := sso.LoadSSOConfig(ctx)
	if ssoErr != nil {
		// SSO 設定読み込みエラー（レガシー形式検出含む）はそのまま返す
		return nil, ssoErr
	}
	if ssoCfg == nil {
		// 非 SSO プロファイル → 元のエラーを返す
		return nil, err
	}

	// SSO セッション期限切れ → JSON イベント出力
	enc := json.NewEncoder(events)
	_ = enc.Encode(jsonout.NewSSOSessionExpiredEvent())
	if loginErr := sso.Login(ctx, events); loginErr != nil {
		return nil, fmt.Errorf("SSO login failed: %w", loginErr)
	}

	// ログイン成功後にリトライ（1回のみ）
	return loadCredentialsOnce(ctx)
}
