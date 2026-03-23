package sso

import (
	"context"
	"fmt"
	"io"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/youyo/awslogin/browse"
)

const defaultExpiresInSeconds int32 = 3600

// Login は SSO OIDC デバイス認証フローを実行してトークンをキャッシュする
func Login(ctx context.Context, events io.Writer) error {
	cfg, err := LoadSSOConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load SSO config: %w", err)
	}
	if cfg == nil {
		return fmt.Errorf("current profile is not configured for SSO")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.SSORegion))
	if err != nil {
		return fmt.Errorf("failed to load AWS config for SSO region: %w", err)
	}

	client := ssooidc.NewFromConfig(awsCfg)

	result, err := RunDeviceAuthFlow(ctx, client, cfg, browse.Start, events)
	if err != nil {
		return err
	}

	// トークンをキャッシュに保存
	cacheKey := CacheKey(cfg)
	cachePath, err := CacheFilePath(cacheKey)
	if err != nil {
		return fmt.Errorf("failed to determine cache path: %w", err)
	}

	expiresIn := result.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = defaultExpiresInSeconds
	}

	token := &CachedToken{
		AccessToken:  result.AccessToken,
		ExpiresAt:    time.Now().Add(time.Duration(expiresIn) * time.Second).UTC().Format(time.RFC3339),
		RefreshToken: result.RefreshToken,
		ClientID:     result.ClientID,
		ClientSecret: result.ClientSecret,
		StartURL:     cfg.StartURL,
		Region:       cfg.SSORegion,
	}

	if err := WriteToken(cachePath, token); err != nil {
		return fmt.Errorf("failed to cache SSO token: %w", err)
	}

	return nil
}
