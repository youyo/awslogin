package awslogin

import (
	"github.com/skratchdot/open-golang/open"
	latest "github.com/tcnksm/go-latest"
)

func VersionCheck(version string) (outdate bool, currentVersion string, err error) {
	githubTag := &latest.GithubTag{
		Owner:      "youyo",
		Repository: "awslogin",
	}
	res, err := latest.Check(githubTag, version)
	if err == nil {
		if res.Outdated {
			outdate = true
			currentVersion = res.Current
			return
		}
	} else {
		return
	}
	return
}

func CheckArgProfileName(r string) bool {
	if r == "" {
		return false
	}
	return true
}

func Browsing(url string) (err error) {
	err = open.Start(url)
	return
}

func BrowsingSpecificApp(url, app string) (err error) {
	err = open.StartWith(url, app)
	return
}
