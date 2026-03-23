package profile

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Profile は AWS プロファイル情報
type Profile struct {
	Name        string `json:"name"`
	Type        string `json:"type"`                    // "sso" | "credentials" | "assume_role"
	SSOStartURL string `json:"sso_start_url,omitempty"`
	Region      string `json:"region,omitempty"`
}

// CurrentSession はカレントシェルの一時認証情報
type CurrentSession struct {
	AWSAccessKeyID  string `json:"aws_access_key_id"`
	HasSessionToken bool   `json:"has_session_token"`
	Region          string `json:"region,omitempty"`
	Source          string `json:"source"`
}

// DefaultConfigPath は ~/.aws/config のデフォルトパスを返す
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".aws", "config")
	}
	return filepath.Join(home, ".aws", "config")
}

// ListProfiles は configPath の AWS config を読み込みプロファイル一覧を返す。
// ファイルが存在しない場合や空ファイルの場合は空スライスと nil エラーを返す。
func ListProfiles(configPath string) ([]Profile, error) {
	f, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Profile{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var profiles []Profile
	var current *Profile
	var inProfileSection bool

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 空行・コメント行はスキップ
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// セクションヘッダー
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// 前のプロファイルを確定
			if current != nil {
				profiles = append(profiles, *current)
				current = nil
			}

			inner := line[1 : len(line)-1]
			name, ok := parseProfileSection(inner)
			if !ok {
				// [sso-session ...] など、プロファイル以外のセクションはスキップ
				inProfileSection = false
				continue
			}

			inProfileSection = true
			current = &Profile{
				Name: name,
				Type: "credentials", // デフォルト; キー解析で上書き
			}
			continue
		}

		// キー=バリュー行
		if !inProfileSection || current == nil {
			continue
		}

		key, value, ok := parseKeyValue(line)
		if !ok {
			continue
		}

		switch key {
		case "sso_start_url":
			current.SSOStartURL = value
			current.Type = "sso"
		case "role_arn":
			current.Type = "assume_role"
		case "region":
			current.Region = value
		}
	}

	// 最後のプロファイルを確定
	if current != nil {
		profiles = append(profiles, *current)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

// parseProfileSection はセクション内部文字列（[]を除いた部分）を解析し、
// プロファイル名と「プロファイルセクションかどうか」を返す。
//
//   - "profile xxx"  → ("xxx", true)
//   - "default"      → ("default", true)
//   - "sso-session x"→ ("", false)
func parseProfileSection(inner string) (string, bool) {
	if inner == "default" {
		return "default", true
	}
	if strings.HasPrefix(inner, "profile ") {
		name := strings.TrimSpace(strings.TrimPrefix(inner, "profile "))
		if name != "" {
			return name, true
		}
	}
	return "", false
}

// parseKeyValue は "key = value" または "key=value" を解析する。
func parseKeyValue(line string) (key, value string, ok bool) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	return key, value, key != ""
}

// DetectCurrentSession は環境変数から現在のセッション情報を検出する。
// AWS_ACCESS_KEY_ID が未設定の場合は nil を返す。
// AWS_SECRET_ACCESS_KEY は絶対に出力しない。
func DetectCurrentSession() *CurrentSession {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		return nil
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}

	return &CurrentSession{
		AWSAccessKeyID:  MaskAccessKeyID(accessKeyID),
		HasSessionToken: os.Getenv("AWS_SESSION_TOKEN") != "",
		Region:          region,
		Source:          "environment",
	}
}

// MaskAccessKeyID はアクセスキー ID をマスクする。
// 4 文字以下はそのまま返し、5 文字以上は先頭 4 文字 + "..." を返す。
func MaskAccessKeyID(id string) string {
	if len(id) <= 4 {
		return id
	}
	return id[:4] + "..."
}
