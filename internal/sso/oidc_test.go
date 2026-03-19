package sso

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	ssooidctypes "github.com/aws/aws-sdk-go-v2/service/ssooidc/types"
)

// mockOIDCClient は OIDCClient の mock 実装
type mockOIDCClient struct {
	registerClientFn         func(ctx context.Context, params *ssooidc.RegisterClientInput, optFns ...func(*ssooidc.Options)) (*ssooidc.RegisterClientOutput, error)
	startDeviceAuthorizationFn func(ctx context.Context, params *ssooidc.StartDeviceAuthorizationInput, optFns ...func(*ssooidc.Options)) (*ssooidc.StartDeviceAuthorizationOutput, error)
	createTokenCalls         int
	createTokenFn            func(callNum int, ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error)
}

func (m *mockOIDCClient) RegisterClient(ctx context.Context, params *ssooidc.RegisterClientInput, optFns ...func(*ssooidc.Options)) (*ssooidc.RegisterClientOutput, error) {
	return m.registerClientFn(ctx, params, optFns...)
}

func (m *mockOIDCClient) StartDeviceAuthorization(ctx context.Context, params *ssooidc.StartDeviceAuthorizationInput, optFns ...func(*ssooidc.Options)) (*ssooidc.StartDeviceAuthorizationOutput, error) {
	return m.startDeviceAuthorizationFn(ctx, params, optFns...)
}

func (m *mockOIDCClient) CreateToken(ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error) {
	m.createTokenCalls++
	return m.createTokenFn(m.createTokenCalls, ctx, params, optFns...)
}

// defaultRegisterClient は成功するデフォルト RegisterClient mock
func defaultRegisterClient() func(ctx context.Context, params *ssooidc.RegisterClientInput, optFns ...func(*ssooidc.Options)) (*ssooidc.RegisterClientOutput, error) {
	return func(ctx context.Context, params *ssooidc.RegisterClientInput, optFns ...func(*ssooidc.Options)) (*ssooidc.RegisterClientOutput, error) {
		return &ssooidc.RegisterClientOutput{
			ClientId:     aws.String("test-client-id"),
			ClientSecret: aws.String("test-client-secret"),
		}, nil
	}
}

// defaultStartDeviceAuth は成功するデフォルト StartDeviceAuthorization mock（Interval=0 で即ポーリング）
func defaultStartDeviceAuth() func(ctx context.Context, params *ssooidc.StartDeviceAuthorizationInput, optFns ...func(*ssooidc.Options)) (*ssooidc.StartDeviceAuthorizationOutput, error) {
	return func(ctx context.Context, params *ssooidc.StartDeviceAuthorizationInput, optFns ...func(*ssooidc.Options)) (*ssooidc.StartDeviceAuthorizationOutput, error) {
		return &ssooidc.StartDeviceAuthorizationOutput{
			DeviceCode:              aws.String("test-device-code"),
			UserCode:                aws.String("ABCD-1234"),
			VerificationUri:         aws.String("https://device.sso.ap-northeast-1.amazonaws.com/"),
			VerificationUriComplete: aws.String("https://device.sso.ap-northeast-1.amazonaws.com/?user_code=ABCD-1234"),
			ExpiresIn:               300,
			Interval:                0, // テスト高速化のため 0 を設定（コード内でデフォルト 5 秒になるが、mock は即応答）
		}, nil
	}
}

func noBrowser(url string) error { return nil }

func TestRunDeviceAuthFlow(t *testing.T) {
	cfg := &SSOConfig{
		StartURL:    "https://example.awsapps.com/start",
		SSORegion:   "ap-northeast-1",
		SessionName: "test-session",
	}

	t.Run("N1: 即成功フロー", func(t *testing.T) {
		client := &mockOIDCClient{
			registerClientFn:         defaultRegisterClient(),
			startDeviceAuthorizationFn: defaultStartDeviceAuth(),
			createTokenFn: func(callNum int, ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error) {
				return &ssooidc.CreateTokenOutput{
					AccessToken:  aws.String("test-access-token"),
					ExpiresIn:    3600,
					RefreshToken: aws.String("test-refresh-token"),
				}, nil
			},
		}

		result, err := RunDeviceAuthFlow(context.Background(), client, cfg, noBrowser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.AccessToken != "test-access-token" {
			t.Errorf("expected access token %q, got %q", "test-access-token", result.AccessToken)
		}
		if result.ClientID != "test-client-id" {
			t.Errorf("expected client id %q, got %q", "test-client-id", result.ClientID)
		}
	})

	t.Run("N2: AuthorizationPending 2回 → 3回目成功", func(t *testing.T) {
		client := &mockOIDCClient{
			registerClientFn:         defaultRegisterClient(),
			startDeviceAuthorizationFn: defaultStartDeviceAuth(),
			createTokenFn: func(callNum int, ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error) {
				if callNum <= 2 {
					return nil, &ssooidctypes.AuthorizationPendingException{
						Message: aws.String("authorization pending"),
					}
				}
				return &ssooidc.CreateTokenOutput{
					AccessToken: aws.String("token-after-pending"),
					ExpiresIn:   3600,
				}, nil
			},
		}

		result, err := RunDeviceAuthFlow(context.Background(), client, cfg, noBrowser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.AccessToken != "token-after-pending" {
			t.Errorf("expected %q, got %q", "token-after-pending", result.AccessToken)
		}
		if client.createTokenCalls != 3 {
			t.Errorf("expected 3 CreateToken calls, got %d", client.createTokenCalls)
		}
	})

	t.Run("N3: SlowDown → interval 増加（フロー継続）", func(t *testing.T) {
		client := &mockOIDCClient{
			registerClientFn:         defaultRegisterClient(),
			startDeviceAuthorizationFn: defaultStartDeviceAuth(),
			createTokenFn: func(callNum int, ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error) {
				if callNum == 1 {
					return nil, &ssooidctypes.SlowDownException{
						Message: aws.String("slow down"),
					}
				}
				return &ssooidc.CreateTokenOutput{
					AccessToken: aws.String("token-after-slowdown"),
					ExpiresIn:   3600,
				}, nil
			},
		}

		result, err := RunDeviceAuthFlow(context.Background(), client, cfg, noBrowser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.AccessToken != "token-after-slowdown" {
			t.Errorf("expected %q, got %q", "token-after-slowdown", result.AccessToken)
		}
	})

	t.Run("E1: RegisterClient 失敗", func(t *testing.T) {
		client := &mockOIDCClient{
			registerClientFn: func(ctx context.Context, params *ssooidc.RegisterClientInput, optFns ...func(*ssooidc.Options)) (*ssooidc.RegisterClientOutput, error) {
				return nil, errors.New("register client failed")
			},
			startDeviceAuthorizationFn: defaultStartDeviceAuth(),
		}

		_, err := RunDeviceAuthFlow(context.Background(), client, cfg, noBrowser)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("E2: StartDeviceAuth 失敗", func(t *testing.T) {
		client := &mockOIDCClient{
			registerClientFn: defaultRegisterClient(),
			startDeviceAuthorizationFn: func(ctx context.Context, params *ssooidc.StartDeviceAuthorizationInput, optFns ...func(*ssooidc.Options)) (*ssooidc.StartDeviceAuthorizationOutput, error) {
				return nil, errors.New("start device auth failed")
			},
		}

		_, err := RunDeviceAuthFlow(context.Background(), client, cfg, noBrowser)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("E3: ExpiredTokenException", func(t *testing.T) {
		client := &mockOIDCClient{
			registerClientFn:         defaultRegisterClient(),
			startDeviceAuthorizationFn: defaultStartDeviceAuth(),
			createTokenFn: func(callNum int, ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error) {
				return nil, &ssooidctypes.ExpiredTokenException{
					Message: aws.String("expired"),
				}
			},
		}

		_, err := RunDeviceAuthFlow(context.Background(), client, cfg, noBrowser)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("X1: context キャンセル", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		callCount := 0
		client := &mockOIDCClient{
			registerClientFn:         defaultRegisterClient(),
			startDeviceAuthorizationFn: defaultStartDeviceAuth(),
			createTokenFn: func(callNum int, ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error) {
				callCount++
				if callCount == 1 {
					cancel() // 1回目の呼び出しでキャンセル
				}
				return nil, &ssooidctypes.AuthorizationPendingException{
					Message: aws.String("authorization pending"),
				}
			},
		}

		_, err := RunDeviceAuthFlow(ctx, client, cfg, noBrowser)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("X2: ブラウザ起動失敗でもフロー継続", func(t *testing.T) {
		failBrowser := func(url string) error {
			return errors.New("browser failed to open")
		}

		client := &mockOIDCClient{
			registerClientFn:         defaultRegisterClient(),
			startDeviceAuthorizationFn: defaultStartDeviceAuth(),
			createTokenFn: func(callNum int, ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error) {
				return &ssooidc.CreateTokenOutput{
					AccessToken: aws.String("token-browser-failed"),
					ExpiresIn:   3600,
				}, nil
			},
		}

		result, err := RunDeviceAuthFlow(context.Background(), client, cfg, failBrowser)
		if err != nil {
			t.Fatalf("expected success even with browser failure, got error: %v", err)
		}
		if result.AccessToken != "token-browser-failed" {
			t.Errorf("expected %q, got %q", "token-browser-failed", result.AccessToken)
		}
	})
}
