package sso

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	ssooidctypes "github.com/aws/aws-sdk-go-v2/service/ssooidc/types"
	"github.com/youyo/awslogin/internal/jsonout"
)

const (
	oidcClientName              = "awslogin"
	oidcClientType              = "public"
	oidcScope                   = "sso:account:access"
	oidcGrantType               = "urn:ietf:params:oauth:grant-type:device_code"
	defaultPollInterval   int32 = 5
)

// OIDCClient は SSO OIDC API のインターフェース（テスト用 mock 対応）
type OIDCClient interface {
	RegisterClient(ctx context.Context, params *ssooidc.RegisterClientInput, optFns ...func(*ssooidc.Options)) (*ssooidc.RegisterClientOutput, error)
	StartDeviceAuthorization(ctx context.Context, params *ssooidc.StartDeviceAuthorizationInput, optFns ...func(*ssooidc.Options)) (*ssooidc.StartDeviceAuthorizationOutput, error)
	CreateToken(ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error)
}

// DeviceAuthResult はデバイス認証フローの結果
type DeviceAuthResult struct {
	AccessToken  string
	ExpiresIn    int32
	RefreshToken string
	ClientID     string
	ClientSecret string
}

// RunDeviceAuthFlow は OIDC デバイス認証フローを実行する
func RunDeviceAuthFlow(ctx context.Context, client OIDCClient, cfg *SSOConfig, openBrowser func(string) error, events io.Writer) (*DeviceAuthResult, error) {
	// 1. RegisterClient
	regOut, err := client.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String(oidcClientName),
		ClientType: aws.String(oidcClientType),
		Scopes:     []string{oidcScope},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register client: %w", err)
	}

	// 2. StartDeviceAuthorization
	authOut, err := client.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     regOut.ClientId,
		ClientSecret: regOut.ClientSecret,
		StartUrl:     aws.String(cfg.StartURL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start device authorization: %w", err)
	}

	// 3. stderr に JSON イベント出力
	enc := json.NewEncoder(events)
	verificationURL := aws.ToString(authOut.VerificationUriComplete)
	verificationURLBase := aws.ToString(authOut.VerificationUri)
	if verificationURL == "" {
		verificationURL = verificationURLBase
	}
	_ = enc.Encode(jsonout.NewSSOAuthRequiredEvent(
		aws.ToString(authOut.UserCode),
		verificationURL,
		verificationURLBase,
	))

	// 4. ブラウザ起動（失敗しても続行）
	if authOut.VerificationUriComplete != nil {
		if err := openBrowser(aws.ToString(authOut.VerificationUriComplete)); err != nil {
			_ = enc.Encode(jsonout.NewBrowserOpenFailedEvent())
		} else {
			_ = enc.Encode(jsonout.NewBrowserOpenedEvent(aws.ToString(authOut.VerificationUriComplete)))
		}
	}

	// 5. ポーリング with タイムアウト
	pollCtx, cancel := context.WithTimeout(ctx, time.Duration(authOut.ExpiresIn)*time.Second)
	defer cancel()

	// AWS API が返す Interval を使う。0 の場合は 5 秒をデフォルトとする。
	// テストでは mock が Interval=0 を返すことで time.After(0) = 即時となり高速化できる。
	interval := authOut.Interval
	if interval < 0 {
		interval = defaultPollInterval
	}

	deviceCode := aws.ToString(authOut.DeviceCode)

	for {
		select {
		case <-pollCtx.Done():
			return nil, fmt.Errorf("device authorization timed out: %w", pollCtx.Err())
		case <-time.After(time.Duration(interval) * time.Second):
		}

		tokenOut, err := client.CreateToken(pollCtx, &ssooidc.CreateTokenInput{
			ClientId:     regOut.ClientId,
			ClientSecret: regOut.ClientSecret,
			DeviceCode:   aws.String(deviceCode),
			GrantType:    aws.String(oidcGrantType),
		})
		if err != nil {
			var authPending *ssooidctypes.AuthorizationPendingException
			var slowDown *ssooidctypes.SlowDownException
			var expired *ssooidctypes.ExpiredTokenException

			if errors.As(err, &authPending) {
				continue
			}
			if errors.As(err, &slowDown) {
				interval += defaultPollInterval
				continue
			}
			if errors.As(err, &expired) {
				return nil, fmt.Errorf("device authorization expired: %w", err)
			}
			return nil, fmt.Errorf("failed to create token: %w", err)
		}

		_ = enc.Encode(jsonout.NewSSOAuthCompleteEvent())

		return &DeviceAuthResult{
			AccessToken:  aws.ToString(tokenOut.AccessToken),
			ExpiresIn:    tokenOut.ExpiresIn,
			RefreshToken: aws.ToString(tokenOut.RefreshToken),
			ClientID:     aws.ToString(regOut.ClientId),
			ClientSecret: aws.ToString(regOut.ClientSecret),
		}, nil
	}
}
