//go:build windows

package browse

import "os/exec"

func openURL(url string) *exec.Cmd {
	return exec.Command("cmd", "/c", "start", "", url)
}
