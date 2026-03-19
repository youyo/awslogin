package sso

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CachedToken は ~/.aws/sso/cache/ に保存される SSO トークン
type CachedToken struct {
	AccessToken  string `json:"accessToken"`
	ExpiresAt    string `json:"expiresAt"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	StartURL     string `json:"startUrl,omitempty"`
	Region       string `json:"region,omitempty"`
}

// CacheFilePath は session name/start URL からキャッシュファイルパスを生成する
// SHA1(key) → 小文字hex → ~/.aws/sso/cache/<hash>.json
func CacheFilePath(key string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	h := sha1.New()
	h.Write([]byte(key))
	hash := strings.ToLower(hex.EncodeToString(h.Sum(nil)))

	return filepath.Join(homeDir, ".aws", "sso", "cache", hash+".json"), nil
}

// WriteToken はトークンを JSON ファイルに書き込む（パーミッション 0600）
// ディレクトリが存在しない場合は 0700 で自動作成
func WriteToken(path string, token *CachedToken) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// ReadToken は JSON ファイルからトークンを読み込む
func ReadToken(path string) (*CachedToken, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var token CachedToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}
