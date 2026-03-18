//go:build linux

package browse

import "os/exec"

func openURL(url string) *exec.Cmd {
	return exec.Command("xdg-open", url)
}
