package sso

import (
	"errors"
	"fmt"
	"testing"

	ssooidctypes "github.com/aws/aws-sdk-go-v2/service/ssooidc/types"
)

func TestIsInvalidGrantError(t *testing.T) {
	t.Run("N6: InvalidGrantException 直接", func(t *testing.T) {
		err := &ssooidctypes.InvalidGrantException{}
		if !IsInvalidGrantError(err) {
			t.Error("expected true, got false")
		}
	})

	t.Run("N7: ラップされた InvalidGrantException", func(t *testing.T) {
		inner := &ssooidctypes.InvalidGrantException{}
		err := fmt.Errorf("wrap: %w", inner)
		if !IsInvalidGrantError(err) {
			t.Error("expected true, got false")
		}
	})

	t.Run("E6: 無関係なエラー", func(t *testing.T) {
		err := errors.New("other error")
		if IsInvalidGrantError(err) {
			t.Error("expected false, got true")
		}
	})

	t.Run("X3: nil", func(t *testing.T) {
		if IsInvalidGrantError(nil) {
			t.Error("expected false, got true")
		}
	})
}
