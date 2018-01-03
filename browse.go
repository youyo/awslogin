package awslogin

import (
	"github.com/skratchdot/open-golang/open"
)

func Browsing(url string) (err error) {
	err = open.Start(url)
	return
}

func BrowsingSpecificApp(url, app string) (err error) {
	err = open.StartWith(url, app)
	return
}
