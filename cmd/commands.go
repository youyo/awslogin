package cmd

import (
	"context"
	"fmt"

	"github.com/youyo/awslogin/browse"
	"github.com/youyo/awslogin/internal/signin"
)

// LoginCmd はデフォルトコマンド。AWS マネジメントコンソールのログイン URL を生成する。
type LoginCmd struct{}

// Run はログイン URL を生成して stdout に出力する
func (c *LoginCmd) Run(globals *Globals) error {
	ctx := context.Background()

	// 1. AWS 認証情報取得
	creds, err := signin.LoadCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
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

	// 3. SigninToken 取得 URL 構築
	requestURL := signin.BuildSigninTokenRequestURL(temporaryCredentials, globals.Duration)

	// 4. SigninToken 取得
	signinToken, err := signin.RequestSigninToken(requestURL)
	if err != nil {
		return fmt.Errorf("failed to request signin token: %w", err)
	}

	// 5. ログイン URL 生成
	signinURL := signin.BuildSigninURL(signinToken, creds.Region)

	// 6. 出力: --open 時はブラウザで開く、それ以外は stdout に出力
	if globals.Open {
		return browse.Start(signinURL)
	}
	fmt.Println(signinURL)

	return nil
}

// VersionCmd はバージョン情報を表示する
type VersionCmd struct{}

// Run はバージョン情報を出力する
func (c *VersionCmd) Run(globals *Globals) error {
	fmt.Printf("awslogin version %s (commit %s, built %s)\n",
		globals.AppVersion, globals.Commit, globals.Date)
	return nil
}
