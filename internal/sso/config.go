package sso

import (
	"context"
	"fmt"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

// SSOConfig は SSO プロファイルの設定情報
type SSOConfig struct {
	StartURL    string
	SSORegion   string
	SessionName string
}

// LoadSSOConfig は現在のプロファイルから SSO 設定を読み込む
// 非 SSO プロファイル → nil, nil（エラーではない）
// レガシー形式（sso_start_url 直接指定） → エラー + 移行ガイド
func LoadSSOConfig(ctx context.Context) (*SSOConfig, error) {
	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "default"
	}

	sharedCfg, err := awsconfig.LoadSharedConfigProfile(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS profile %q: %w", profile, err)
	}

	// 新形式: sso-session セクションあり
	if sharedCfg.SSOSession != nil {
		return &SSOConfig{
			StartURL:    sharedCfg.SSOSession.SSOStartURL,
			SSORegion:   sharedCfg.SSOSession.SSORegion,
			SessionName: sharedCfg.SSOSession.Name,
		}, nil
	}

	// レガシー形式: sso_start_url 直接指定
	if sharedCfg.SSOStartURL != "" {
		return nil, fmt.Errorf(
			"legacy SSO configuration detected for profile %q.\n"+
				"Please migrate to the sso-session format:\n\n"+
				"  [profile %s]\n"+
				"  sso_session = my-sso\n\n"+
				"  [sso-session my-sso]\n"+
				"  sso_start_url = %s\n"+
				"  sso_region = %s\n",
			profile, profile, sharedCfg.SSOStartURL, sharedCfg.SSORegion,
		)
	}

	// 非 SSO プロファイル
	return nil, nil
}

// CacheKey は SSO 設定からキャッシュキーを返す
func CacheKey(cfg *SSOConfig) string {
	return cfg.SessionName
}
