//go:build darwin || linux

package browse

import (
	"runtime"
	"testing"
)

func TestOpenURL(t *testing.T) {
	url := "https://signin.aws.amazon.com/federation?Action=login&Destination=https%3A%2F%2Fconsole.aws.amazon.com%2F"
	cmd := openURL(url)

	if len(cmd.Args) < 2 {
		t.Fatalf("expected at least 2 args, got %d", len(cmd.Args))
	}

	// 最後の引数が URL であること
	lastArg := cmd.Args[len(cmd.Args)-1]
	if lastArg != url {
		t.Errorf("expected last arg to be URL, got %s", lastArg)
	}

	// プラットフォーム別のコマンド名を検証
	switch runtime.GOOS {
	case "darwin":
		if cmd.Args[0] != "open" {
			t.Errorf("expected command 'open' on darwin, got %s", cmd.Args[0])
		}
	case "linux":
		if cmd.Args[0] != "xdg-open" {
			t.Errorf("expected command 'xdg-open' on linux, got %s", cmd.Args[0])
		}
	}
}

func TestOpenURLArgsCount(t *testing.T) {
	cmd := openURL("https://example.com")

	switch runtime.GOOS {
	case "darwin":
		// open <url> → 2 args
		if len(cmd.Args) != 2 {
			t.Errorf("expected 2 args on darwin, got %d: %v", len(cmd.Args), cmd.Args)
		}
	case "linux":
		// xdg-open <url> → 2 args
		if len(cmd.Args) != 2 {
			t.Errorf("expected 2 args on linux, got %d: %v", len(cmd.Args), cmd.Args)
		}
	}
}
