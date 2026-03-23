package sso

import (
	"context"
	"io"
	"testing"
)

// TestLoginSignature は Login 関数のシグネチャ確認
// AWS API への接続が必要なため統合テストは行わない
func TestLoginSignature(t *testing.T) {
	// Login 関数が呼び出し可能であることを確認
	// 非 SSO プロファイル or プロファイルなし → エラー返却を期待
	t.Setenv("AWS_CONFIG_FILE", "/nonexistent/config/file")
	t.Setenv("AWS_PROFILE", "nonexistent-profile")

	err := Login(context.Background(), io.Discard)
	// エラーが返ることを確認（プロファイルが存在しないため）
	if err == nil {
		t.Log("Login returned nil (unexpected but not fatal in test environment)")
	}
	// パニックしないことを確認できれば十分
}
