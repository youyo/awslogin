package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompletionCmd_Zsh(t *testing.T) {
	var buf bytes.Buffer
	c := &CompletionCmd{
		Shell:  "zsh",
		Writer: &buf,
	}

	if err := c.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "compdef") {
		t.Error("zsh completion script should contain 'compdef'")
	}
	if !strings.Contains(output, "_awslogin") {
		t.Error("zsh completion script should contain '_awslogin'")
	}
	if !strings.Contains(output, "awslogin") {
		t.Error("zsh completion script should reference 'awslogin'")
	}
}

func TestCompletionCmd_Bash(t *testing.T) {
	var buf bytes.Buffer
	c := &CompletionCmd{
		Shell:  "bash",
		Writer: &buf,
	}

	if err := c.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "complete -F") {
		t.Error("bash completion script should contain 'complete -F'")
	}
	if !strings.Contains(output, "awslogin") {
		t.Error("bash completion script should reference 'awslogin'")
	}
}

func TestCompletionCmd_DefaultWriter(t *testing.T) {
	// Writer が nil の場合でもパニックしないことを確認
	c := &CompletionCmd{
		Shell: "zsh",
	}

	// os.Stdout に書き込まれるのでエラーがないことだけ確認
	if err := c.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
