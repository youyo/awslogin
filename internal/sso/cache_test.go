package sso

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
)

func TestCacheFilePath(t *testing.T) {
	t.Run("N5: 一貫した hex 文字列", func(t *testing.T) {
		path1, err := CacheFilePath("my-sso-session")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		path2, err := CacheFilePath("my-sso-session")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path1 != path2 {
			t.Errorf("expected consistent path, got %q and %q", path1, path2)
		}
		// パスが .json で終わる
		if !strings.HasSuffix(path1, ".json") {
			t.Errorf("expected .json suffix, got %q", path1)
		}
		// ファイル名部分が小文字 hex のみ
		base := filepath.Base(path1)
		name := strings.TrimSuffix(base, ".json")
		for _, c := range name {
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
				t.Errorf("expected lowercase hex filename, got %q", name)
				break
			}
		}
	})

	t.Run("C2: Go SDK の ssocreds.StandardCachedTokenFilepath と一致", func(t *testing.T) {
		key := "my-sso-session"
		got, err := CacheFilePath(key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want, err := ssocreds.StandardCachedTokenFilepath(key)
		if err != nil {
			t.Fatalf("ssocreds error: %v", err)
		}
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}

func TestWriteReadToken(t *testing.T) {
	tmpDir := os.Getenv("TMPDIR")
	if tmpDir == "" {
		tmpDir = "/private/tmp/claude-501/"
	}

	t.Run("N4: WriteToken → ReadToken 往復テスト", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "sso-cache-test-n4")
		defer func() { _ = os.RemoveAll(dir) }()

		path := filepath.Join(dir, "token.json")
		token := &CachedToken{
			AccessToken:  "test-access-token",
			ExpiresAt:    "2024-01-01T00:00:00Z",
			RefreshToken: "test-refresh-token",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			StartURL:     "https://example.awsapps.com/start",
			Region:       "ap-northeast-1",
		}

		if err := WriteToken(path, token); err != nil {
			t.Fatalf("WriteToken failed: %v", err)
		}

		got, err := ReadToken(path)
		if err != nil {
			t.Fatalf("ReadToken failed: %v", err)
		}

		if got.AccessToken != token.AccessToken {
			t.Errorf("AccessToken: expected %q, got %q", token.AccessToken, got.AccessToken)
		}
		if got.ExpiresAt != token.ExpiresAt {
			t.Errorf("ExpiresAt: expected %q, got %q", token.ExpiresAt, got.ExpiresAt)
		}
		if got.RefreshToken != token.RefreshToken {
			t.Errorf("RefreshToken: expected %q, got %q", token.RefreshToken, got.RefreshToken)
		}
		if got.ClientID != token.ClientID {
			t.Errorf("ClientID: expected %q, got %q", token.ClientID, got.ClientID)
		}
		if got.ClientSecret != token.ClientSecret {
			t.Errorf("ClientSecret: expected %q, got %q", token.ClientSecret, got.ClientSecret)
		}
		if got.StartURL != token.StartURL {
			t.Errorf("StartURL: expected %q, got %q", token.StartURL, got.StartURL)
		}
		if got.Region != token.Region {
			t.Errorf("Region: expected %q, got %q", token.Region, got.Region)
		}
	})

	t.Run("E4: ReadToken 存在しないファイル", func(t *testing.T) {
		_, err := ReadToken(filepath.Join(tmpDir, "nonexistent-xyz/token.json"))
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("E5: ReadToken 壊れた JSON", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "sso-cache-test-e5")
		defer func() { _ = os.RemoveAll(dir) }()
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, "broken.json")
		if err := os.WriteFile(path, []byte("{not valid json}"), 0600); err != nil {
			t.Fatal(err)
		}
		_, err := ReadToken(path)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("X4: WriteToken でディレクトリ自動作成", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "sso-cache-test-x4", "nested", "dir")
		defer func() { _ = os.RemoveAll(filepath.Join(tmpDir, "sso-cache-test-x4")) }()

		path := filepath.Join(dir, "token.json")
		token := &CachedToken{
			AccessToken: "auto-create-test",
			ExpiresAt:   "2024-01-01T00:00:00Z",
		}
		if err := WriteToken(path, token); err != nil {
			t.Fatalf("WriteToken failed: %v", err)
		}
		if _, err := os.Stat(path); err != nil {
			t.Errorf("file not created: %v", err)
		}
	})

	t.Run("C1: JSON に全フィールドが含まれる", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "sso-cache-test-c1")
		defer func() { _ = os.RemoveAll(dir) }()
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}

		path := filepath.Join(dir, "token.json")
		token := &CachedToken{
			AccessToken:  "access",
			ExpiresAt:    "2024-01-01T00:00:00Z",
			StartURL:     "https://example.awsapps.com/start",
			Region:       "us-east-1",
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		}
		if err := WriteToken(path, token); err != nil {
			t.Fatalf("WriteToken failed: %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatal(err)
		}

		for _, field := range []string{"startUrl", "region", "clientId", "clientSecret"} {
			if _, ok := m[field]; !ok {
				t.Errorf("expected JSON field %q not found", field)
			}
		}
	})

	t.Run("C3: expiresAt が RFC3339 形式で往復する", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "sso-cache-test-c3")
		defer func() { _ = os.RemoveAll(dir) }()
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}

		now := time.Now().UTC().Truncate(time.Second)
		expiresAt := now.Format(time.RFC3339)

		path := filepath.Join(dir, "token.json")
		token := &CachedToken{
			AccessToken: "access",
			ExpiresAt:   expiresAt,
		}
		if err := WriteToken(path, token); err != nil {
			t.Fatalf("WriteToken failed: %v", err)
		}

		got, err := ReadToken(path)
		if err != nil {
			t.Fatalf("ReadToken failed: %v", err)
		}

		if got.ExpiresAt != expiresAt {
			t.Errorf("ExpiresAt: expected %q, got %q", expiresAt, got.ExpiresAt)
		}

		// RFC3339 でパース可能であることを確認
		if _, err := time.Parse(time.RFC3339, got.ExpiresAt); err != nil {
			t.Errorf("ExpiresAt not RFC3339: %v", err)
		}
	})
}

func TestCacheFilePermissions(t *testing.T) {
	tmpDir := os.Getenv("TMPDIR")
	if tmpDir == "" {
		tmpDir = "/private/tmp/claude-501/"
	}

	t.Run("ファイルパーミッション 0600", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "sso-cache-perm-test")
		defer func() { _ = os.RemoveAll(dir) }()
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}

		path := filepath.Join(dir, "token.json")
		token := &CachedToken{
			AccessToken: "perm-test",
			ExpiresAt:   "2024-01-01T00:00:00Z",
		}
		if err := WriteToken(path, token); err != nil {
			t.Fatalf("WriteToken failed: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0600 {
			t.Errorf("expected 0600, got %o", info.Mode().Perm())
		}
	})
}
