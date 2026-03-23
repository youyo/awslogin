package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/youyo/awslogin/internal/profile"
)

// --- ListProfiles ---

// P1: SSO 設定を含む config → type="sso", sso_start_url 設定
func TestListProfiles_SSOProfiles(t *testing.T) {
	content := `[profile dev]
sso_start_url = https://example.awsapps.com/start
sso_account_id = 123456789012
sso_role_name = Developer
region = ap-northeast-1

[profile prod]
sso_start_url = https://example.awsapps.com/start
sso_account_id = 999999999999
sso_role_name = ReadOnly
region = us-east-1
`
	configPath := writeTempConfig(t, content)

	profiles, err := profile.ListProfiles(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}

	dev := findProfile(profiles, "dev")
	if dev == nil {
		t.Fatal("profile 'dev' not found")
	}
	if dev.Type != "sso" {
		t.Errorf("expected type 'sso', got %q", dev.Type)
	}
	if dev.SSOStartURL != "https://example.awsapps.com/start" {
		t.Errorf("unexpected SSOStartURL: %q", dev.SSOStartURL)
	}
	if dev.Region != "ap-northeast-1" {
		t.Errorf("unexpected Region: %q", dev.Region)
	}
}

// P2: SSO + credentials + assume_role 混在
func TestListProfiles_MixedProfiles(t *testing.T) {
	content := `[profile sso-user]
sso_start_url = https://example.awsapps.com/start
region = ap-northeast-1

[profile creds-user]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG

[profile role-user]
role_arn = arn:aws:iam::123456789012:role/MyRole
source_profile = creds-user
region = eu-west-1
`
	configPath := writeTempConfig(t, content)

	profiles, err := profile.ListProfiles(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(profiles))
	}

	checkType(t, profiles, "sso-user", "sso")
	checkType(t, profiles, "creds-user", "credentials")
	checkType(t, profiles, "role-user", "assume_role")
}

// P3: 空ファイル → 空スライス
func TestListProfiles_EmptyConfig(t *testing.T) {
	configPath := writeTempConfig(t, "")

	profiles, err := profile.ListProfiles(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 0 {
		t.Fatalf("expected 0 profiles, got %d", len(profiles))
	}
}

// P4: 存在しないパス → 空スライス、エラーなし
func TestListProfiles_NoConfigFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "nonexistent_config")

	profiles, err := profile.ListProfiles(configPath)
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if len(profiles) != 0 {
		t.Fatalf("expected 0 profiles, got %d", len(profiles))
	}
}

// P8: [default] セクション → name="default"
func TestListProfiles_DefaultProfile(t *testing.T) {
	content := `[default]
region = us-east-1
aws_access_key_id = AKIAIOSFODNN7EXAMPLE

[sso-session my-sso]
sso_start_url = https://example.awsapps.com/start
sso_region = ap-northeast-1
`
	configPath := writeTempConfig(t, content)

	profiles, err := profile.ListProfiles(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile (sso-session should be skipped), got %d: %+v", len(profiles), profiles)
	}

	def := findProfile(profiles, "default")
	if def == nil {
		t.Fatal("profile 'default' not found")
	}
	if def.Type != "credentials" {
		t.Errorf("expected type 'credentials', got %q", def.Type)
	}
	if def.Region != "us-east-1" {
		t.Errorf("unexpected Region: %q", def.Region)
	}
}

// --- DetectCurrentSession ---

// P5: AWS_ACCESS_KEY_ID 設定時 → CurrentSession 返却
func TestDetectCurrentSession_WithEnvVars(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "ASIAIOSFODNN7EXAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	t.Setenv("AWS_SESSION_TOKEN", "FQoGZXIvYXdzEJr//////////wEaDummy")
	t.Setenv("AWS_REGION", "ap-northeast-1")
	t.Setenv("AWS_DEFAULT_REGION", "")

	sess := profile.DetectCurrentSession()
	if sess == nil {
		t.Fatal("expected non-nil CurrentSession")
	}
	if sess.AWSAccessKeyID != "ASIA..." {
		t.Errorf("expected masked key 'ASIA...', got %q", sess.AWSAccessKeyID)
	}
	if !sess.HasSessionToken {
		t.Error("expected HasSessionToken to be true")
	}
	if sess.Region != "ap-northeast-1" {
		t.Errorf("expected region 'ap-northeast-1', got %q", sess.Region)
	}
	if sess.Source != "environment" {
		t.Errorf("expected source 'environment', got %q", sess.Source)
	}
}

// P6: AWS_ACCESS_KEY_ID 未設定時 → nil
func TestDetectCurrentSession_NoEnvVars(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("AWS_SESSION_TOKEN", "")

	sess := profile.DetectCurrentSession()
	if sess != nil {
		t.Fatalf("expected nil, got %+v", sess)
	}
}

// DetectCurrentSession: AWS_DEFAULT_REGION フォールバック
func TestDetectCurrentSession_DefaultRegionFallback(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
	t.Setenv("AWS_REGION", "")
	t.Setenv("AWS_DEFAULT_REGION", "eu-west-1")
	t.Setenv("AWS_SESSION_TOKEN", "")

	sess := profile.DetectCurrentSession()
	if sess == nil {
		t.Fatal("expected non-nil CurrentSession")
	}
	if sess.Region != "eu-west-1" {
		t.Errorf("expected region 'eu-west-1', got %q", sess.Region)
	}
}

// --- MaskAccessKeyID ---

// P7: マスク処理のテスト
func TestMaskAccessKeyID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"ABC", "ABC"},
		{"ABCD", "ABCD"},
		{"ABCDE", "ABCD..."},
		{"ASIAIOSFODNN7EXAMPLE", "ASIA..."},
		{"AKIAIOSFODNN7EXAMPLE", "AKIA..."},
	}

	for _, tt := range tests {
		got := profile.MaskAccessKeyID(tt.input)
		if got != tt.expected {
			t.Errorf("MaskAccessKeyID(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// --- helpers ---

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config")
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return configPath
}

func findProfile(profiles []profile.Profile, name string) *profile.Profile {
	for i := range profiles {
		if profiles[i].Name == name {
			return &profiles[i]
		}
	}
	return nil
}

func checkType(t *testing.T, profiles []profile.Profile, name, expectedType string) {
	t.Helper()
	p := findProfile(profiles, name)
	if p == nil {
		t.Errorf("profile %q not found", name)
		return
	}
	if p.Type != expectedType {
		t.Errorf("profile %q: expected type %q, got %q", name, expectedType, p.Type)
	}
}
