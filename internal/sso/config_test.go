package sso

import (
	"context"
	"testing"
)

func TestCacheKey(t *testing.T) {
	t.Run("SessionName をそのまま返す", func(t *testing.T) {
		cfg := &SSOConfig{
			StartURL:    "https://example.awsapps.com/start",
			SSORegion:   "ap-northeast-1",
			SessionName: "my-sso-session",
		}
		got := CacheKey(cfg)
		if got != "my-sso-session" {
			t.Errorf("expected %q, got %q", "my-sso-session", got)
		}
	})

	t.Run("空の SessionName", func(t *testing.T) {
		cfg := &SSOConfig{}
		got := CacheKey(cfg)
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})
}

func TestLoadSSOConfigSignature(t *testing.T) {
	// LoadSSOConfig のシグネチャ確認
	// 実際の AWS プロファイル読み込みは統合テストレベルのため、
	// ここでは関数が呼び出し可能であることを確認する。
	// AWS_CONFIG_FILE を存在しないファイルにすることで、エラーパスを検証する。
	t.Setenv("AWS_CONFIG_FILE", "/nonexistent/config/file")
	t.Setenv("AWS_PROFILE", "nonexistent-profile")

	cfg, err := LoadSSOConfig(context.Background())
	// エラーが返るか、nil, nil が返るかはプロファイル設定次第
	// 少なくともパニックしないことを確認
	_ = cfg
	_ = err
}
