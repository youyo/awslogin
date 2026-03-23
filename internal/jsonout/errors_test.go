package jsonout_test

import (
	"errors"
	"io"
	"testing"

	"github.com/youyo/awslogin/internal/jsonout"
)

// --- J7: TestAppError_Unwrap ---

func TestAppError_Unwrap(t *testing.T) {
	appErr := &jsonout.AppError{
		Code:    jsonout.ErrInternal,
		Message: "wrapped io.EOF",
		Cause:   io.EOF,
	}

	if !errors.Is(appErr, io.EOF) {
		t.Errorf("expected errors.Is(appErr, io.EOF) == true, but got false")
	}
}

func TestAppError_Error(t *testing.T) {
	appErr := &jsonout.AppError{
		Code:    jsonout.ErrInvalidArgs,
		Message: "invalid argument",
	}

	if appErr.Error() != "invalid argument" {
		t.Errorf("expected Error() %q, got %q", "invalid argument", appErr.Error())
	}
}

func TestAllErrorCodes(t *testing.T) {
	codes := []string{
		jsonout.ErrInvalidArgs,
		jsonout.ErrConfigLoadFailed,
		jsonout.ErrSSOSessionExpired,
		jsonout.ErrSSOLoginFailed,
		jsonout.ErrSSOLegacyConfig,
		jsonout.ErrOIDCRegisterFailed,
		jsonout.ErrOIDCDeviceAuthFailed,
		jsonout.ErrOIDCTokenExpired,
		jsonout.ErrFederationAPIFailed,
		jsonout.ErrBrowserOpenFailed,
		jsonout.ErrInternal,
	}

	seen := make(map[string]bool)
	for _, code := range codes {
		if code == "" {
			t.Errorf("found empty error code")
		}
		if seen[code] {
			t.Errorf("duplicate error code: %q", code)
		}
		seen[code] = true
	}
}
