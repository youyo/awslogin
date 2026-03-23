package jsonout_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/youyo/awslogin/internal/jsonout"
)

// --- テスト用型 ---

type LoginResult struct {
	URL string `json:"url"`
}

type VersionResult struct {
	Version string `json:"version"`
}

// --- J1: TestWriteResult_Login ---

func TestWriteResult_Login(t *testing.T) {
	var stdout, stderr bytes.Buffer
	w := jsonout.New(&stdout, &stderr)

	result := LoginResult{URL: "https://example.com/login"}
	if err := w.WriteResult(result); err != nil {
		t.Fatalf("WriteResult failed: %v", err)
	}

	var got struct {
		Result LoginResult `json:"result"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal stdout: %v\nraw: %s", err, stdout.String())
	}

	if got.Result.URL != result.URL {
		t.Errorf("expected URL %q, got %q", result.URL, got.Result.URL)
	}
	if stderr.Len() != 0 {
		t.Errorf("expected stderr to be empty, got %q", stderr.String())
	}
}

// --- J2: TestWriteResult_Version ---

func TestWriteResult_Version(t *testing.T) {
	var stdout, stderr bytes.Buffer
	w := jsonout.New(&stdout, &stderr)

	result := VersionResult{Version: "3.0.0"}
	if err := w.WriteResult(result); err != nil {
		t.Fatalf("WriteResult failed: %v", err)
	}

	var got struct {
		Result VersionResult `json:"result"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal stdout: %v\nraw: %s", err, stdout.String())
	}

	if got.Result.Version != result.Version {
		t.Errorf("expected Version %q, got %q", result.Version, got.Result.Version)
	}
}

// --- J3: TestWriteError_AppError ---

func TestWriteError_AppError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	w := jsonout.New(&stdout, &stderr)

	appErr := &jsonout.AppError{
		Code:    jsonout.ErrSSOSessionExpired,
		Message: "session expired",
		Details: "token expired at 2025-01-01",
	}
	if err := w.WriteError(appErr); err != nil {
		t.Fatalf("WriteError failed: %v", err)
	}

	var got struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal stdout: %v\nraw: %s", err, stdout.String())
	}

	if got.Error.Code != jsonout.ErrSSOSessionExpired {
		t.Errorf("expected code %q, got %q", jsonout.ErrSSOSessionExpired, got.Error.Code)
	}
	if got.Error.Message != appErr.Message {
		t.Errorf("expected message %q, got %q", appErr.Message, got.Error.Message)
	}
	if got.Error.Details != appErr.Details {
		t.Errorf("expected details %q, got %q", appErr.Details, got.Error.Details)
	}
}

// --- J4: TestWriteError_GenericError ---

func TestWriteError_GenericError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	w := jsonout.New(&stdout, &stderr)

	genericErr := errors.New("something went wrong")
	if err := w.WriteError(genericErr); err != nil {
		t.Fatalf("WriteError failed: %v", err)
	}

	var got struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal stdout: %v\nraw: %s", err, stdout.String())
	}

	if got.Error.Code != jsonout.ErrInternal {
		t.Errorf("expected code %q, got %q", jsonout.ErrInternal, got.Error.Code)
	}
	if got.Error.Message != genericErr.Error() {
		t.Errorf("expected message %q, got %q", genericErr.Error(), got.Error.Message)
	}
}

// --- J5: TestWriteError_WrappedAppError ---

func TestWriteError_WrappedAppError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	w := jsonout.New(&stdout, &stderr)

	appErr := &jsonout.AppError{
		Code:    jsonout.ErrSSOLoginFailed,
		Message: "SSO login failed",
	}
	wrappedErr := fmt.Errorf("outer: %w", appErr)
	if err := w.WriteError(wrappedErr); err != nil {
		t.Fatalf("WriteError failed: %v", err)
	}

	var got struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal stdout: %v\nraw: %s", err, stdout.String())
	}

	if got.Error.Code != jsonout.ErrSSOLoginFailed {
		t.Errorf("expected code %q, got %q", jsonout.ErrSSOLoginFailed, got.Error.Code)
	}
}

// --- J6: TestWriteEvent ---

func TestWriteEvent(t *testing.T) {
	var stdout, stderr bytes.Buffer
	w := jsonout.New(&stdout, &stderr)

	event := jsonout.NewSSOAuthRequiredEvent("ABCD-1234", "https://device.example.com/activate?user_code=ABCD-1234", "https://device.example.com/activate")
	if err := w.WriteEvent(event); err != nil {
		t.Fatalf("WriteEvent failed: %v", err)
	}

	if stdout.Len() != 0 {
		t.Errorf("expected stdout to be empty, got %q", stdout.String())
	}

	var got jsonout.SSOAuthRequiredEvent
	if err := json.Unmarshal(stderr.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal stderr: %v\nraw: %s", err, stderr.String())
	}

	if got.Type != "sso_auth_required" {
		t.Errorf("expected type %q, got %q", "sso_auth_required", got.Type)
	}
	if got.VerificationCode != "ABCD-1234" {
		t.Errorf("expected verification_code %q, got %q", "ABCD-1234", got.VerificationCode)
	}
	if got.VerificationURL != "https://device.example.com/activate?user_code=ABCD-1234" {
		t.Errorf("expected verification_url %q, got %q", "https://device.example.com/activate?user_code=ABCD-1234", got.VerificationURL)
	}
	if got.VerificationURLBase != "https://device.example.com/activate" {
		t.Errorf("expected verification_url_base %q, got %q", "https://device.example.com/activate", got.VerificationURLBase)
	}
}
