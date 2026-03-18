package signin

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
)

// AWSCredentials は AWS 認証情報とリージョン情報を保持する
type AWSCredentials struct {
	AccessKeyID    string
	SecretAccessKey string
	SessionToken   string
	Region         string
}

// LoadCredentials は AWS SDK v2 で認証情報とリージョンを取得する。
// AWS_PROFILE 環境変数、~/.aws/credentials、IAM ロール等、SDK の config loader が自動解決する。
func LoadCredentials(ctx context.Context) (*AWSCredentials, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve AWS credentials: %w", err)
	}

	return &AWSCredentials{
		AccessKeyID:    creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:   creds.SessionToken,
		Region:         ResolveRegion(cfg.Region),
	}, nil
}
