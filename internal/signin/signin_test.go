package signin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestBuildTemporaryCredentials(t *testing.T) {
	got, err := BuildTemporaryCredentials("AKID", "SECRET", "TOKEN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var creds TemporaryCredentials
	if err := json.Unmarshal([]byte(got), &creds); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if creds.SessionID != "AKID" {
		t.Errorf("sessionId = %q, want %q", creds.SessionID, "AKID")
	}
	if creds.SessionKey != "SECRET" {
		t.Errorf("sessionKey = %q, want %q", creds.SessionKey, "SECRET")
	}
	if creds.SessionToken != "TOKEN" {
		t.Errorf("sessionToken = %q, want %q", creds.SessionToken, "TOKEN")
	}
}

func TestResolveRegion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "us-east-1"},
		{"ap-northeast-1", "ap-northeast-1"},
		{"eu-west-1", "eu-west-1"},
		{"us-east-1", "us-east-1"},
	}
	for _, tt := range tests {
		if got := ResolveRegion(tt.input); got != tt.want {
			t.Errorf("ResolveRegion(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildSigninTokenRequestURL(t *testing.T) {
	u := BuildSigninTokenRequestURL(`{"sessionId":"x"}`, 3600)
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("invalid URL: %v", err)
	}

	q := parsed.Query()
	if q.Get("Action") != "getSigninToken" {
		t.Errorf("Action = %q, want getSigninToken", q.Get("Action"))
	}
	if q.Get("SessionDuration") != "3600" {
		t.Errorf("SessionDuration = %q, want 3600", q.Get("SessionDuration"))
	}
	if q.Get("Session") != `{"sessionId":"x"}` {
		t.Errorf("Session = %q, want credentials JSON", q.Get("Session"))
	}
	if q.Get("SessionType") != "json" {
		t.Errorf("SessionType = %q, want json", q.Get("SessionType"))
	}
}

func TestBuildSigninURL_APNortheast1(t *testing.T) {
	u := BuildSigninURL("mytoken", "ap-northeast-1")
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("invalid URL: %v", err)
	}

	q := parsed.Query()
	if q.Get("Action") != "login" {
		t.Errorf("Action = %q, want login", q.Get("Action"))
	}
	dest := q.Get("Destination")
	if dest != "https://ap-northeast-1.console.aws.amazon.com/" {
		t.Errorf("Destination = %q, want https://ap-northeast-1.console.aws.amazon.com/", dest)
	}
	if q.Get("SigninToken") != "mytoken" {
		t.Errorf("SigninToken = %q, want mytoken", q.Get("SigninToken"))
	}
}

func TestBuildSigninURL_USEast1(t *testing.T) {
	u := BuildSigninURL("mytoken", "us-east-1")
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("invalid URL: %v", err)
	}

	dest := parsed.Query().Get("Destination")
	if dest != "https://console.aws.amazon.com/" {
		t.Errorf("Destination = %q, want https://console.aws.amazon.com/ (no region prefix)", dest)
	}
}

func TestBuildSigninURL_EmptyRegionFallback(t *testing.T) {
	u := BuildSigninURL("mytoken", "")
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("invalid URL: %v", err)
	}

	dest := parsed.Query().Get("Destination")
	if dest != "https://console.aws.amazon.com/" {
		t.Errorf("empty region should fallback to us-east-1, got Destination = %q", dest)
	}
}

func TestRequestSigninToken_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"SigninToken":"test-token-123"}`))
	}))
	defer ts.Close()

	token, err := RequestSigninToken(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "test-token-123" {
		t.Errorf("got %q, want test-token-123", token)
	}
}

func TestRequestSigninToken_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not json`))
	}))
	defer ts.Close()

	_, err := RequestSigninToken(context.Background(), ts.URL)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestRequestSigninToken_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"Error":"AccessDenied"}`))
	}))
	defer ts.Close()

	_, err := RequestSigninToken(context.Background(), ts.URL)
	if err == nil {
		t.Error("expected error for HTTP 403, got nil")
	}
}
