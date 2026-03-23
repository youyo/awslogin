package cmd

import (
	"context"
	"fmt"

	"github.com/youyo/awslogin/browse"
	"github.com/youyo/awslogin/internal/jsonout"
	"github.com/youyo/awslogin/internal/signin"
)

// LoginCmd はデフォルトコマンド。AWS マネジメントコンソールのログイン URL を生成する。
type LoginCmd struct{}

// LoginResult は login コマンドの JSON 出力
type LoginResult struct {
	URL             string `json:"url"`
	Region          string `json:"region"`
	OpenedInBrowser bool   `json:"opened_in_browser"`
}

// Run はログイン URL を生成して JSON で stdout に出力する
func (c *LoginCmd) Run(globals *Globals, out *jsonout.Writer) error {
	ctx := context.Background()

	// 1. AWS 認証情報取得（events は out の stderr）
	creds, err := signin.LoadCredentials(ctx, out.Stderr())
	if err != nil {
		return err
	}

	// 2. 一時認証情報 JSON 化
	temporaryCredentials, err := signin.BuildTemporaryCredentials(
		creds.AccessKeyID,
		creds.SecretAccessKey,
		creds.SessionToken,
	)
	if err != nil {
		return fmt.Errorf("failed to build temporary credentials: %w", err)
	}

	// 3. SigninToken 取得
	requestURL := signin.BuildSigninTokenRequestURL(temporaryCredentials, globals.Duration)
	signinToken, err := signin.RequestSigninToken(ctx, requestURL)
	if err != nil {
		return err
	}

	// 4. ログイン URL 生成
	signinURL := signin.BuildSigninURL(signinToken, creds.Region)

	// 5. --open 時はブラウザで開く
	openedInBrowser := false
	if globals.Open {
		if err := browse.Start(signinURL); err != nil {
			// ブラウザ失敗は非致命的。URL は返す
			_ = out.WriteEvent(jsonout.NewBrowserOpenFailedEvent())
		} else {
			openedInBrowser = true
		}
	}

	return out.WriteResult(LoginResult{
		URL:             signinURL,
		Region:          creds.Region,
		OpenedInBrowser: openedInBrowser,
	})
}

// VersionCmd はバージョン情報を表示する
type VersionCmd struct{}

// VersionResult は version コマンドの JSON 出力
type VersionResult struct {
	Version string `json:"version"`
}

// Run はバージョン情報を JSON で出力する
func (c *VersionCmd) Run(globals *Globals, out *jsonout.Writer) error {
	return out.WriteResult(VersionResult{Version: globals.AppVersion})
}
