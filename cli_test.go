package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_listFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./awslogin -list", " ")

	if _, exist := os.Stat(filepath.Join(os.Getenv("HOME"), ".aws/config")); exist != nil {
		p := filepath.Join(os.Getenv("HOME"), ".aws")
		f := filepath.Join(os.Getenv("HOME"), ".aws/config")
		_ = os.Mkdir(p, 0700)
		c, _ := os.Create(f)
		c.Close()
	}

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %v to eq %v", status, ExitCodeOK)
	}

	/*
		expected := fmt.Sprintf("")
		if !strings.Contains(errStream.String(), expected) {
			t.Errorf("expected %q to eq %q", errStream.String(), expected)
		}
	*/
}
