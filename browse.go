package awslogin

import (
	"github.com/skratchdot/open-golang/open"
)

func Browsing(url string) error {
	return open.Start(url)
}

func BrowsingSpecificApp(url, app string) error {
	return open.StartWith(url, app)
}
