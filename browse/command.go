// +build !darwin

package browse

import "os/exec"

func open(url string) *exec.Cmd {
	return exec.Command("xdg-open", url)
}

func openWith(url string, app string) *exec.Cmd {
	return exec.Command(app, url)
}
