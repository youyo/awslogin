//go:build windows

package browse

import (
	"testing"
)

func TestOpenURL(t *testing.T) {
	url := "https://signin.aws.amazon.com/federation?Action=login&Destination=https%3A%2F%2Fconsole.aws.amazon.com%2F"
	cmd := openURL(url)

	if len(cmd.Args) < 3 {
		t.Fatalf("expected at least 3 args, got %d: %v", len(cmd.Args), cmd.Args)
	}

	// コマンドが cmd であること
	if cmd.Args[0] != "cmd" {
		t.Errorf("expected command 'cmd', got %s", cmd.Args[0])
	}

	// 最後の引数が URL であること
	lastArg := cmd.Args[len(cmd.Args)-1]
	if lastArg != url {
		t.Errorf("expected last arg to be URL, got %s", lastArg)
	}
}
