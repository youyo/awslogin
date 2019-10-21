package browse

func Start(url string) (err error) {
	err = open(url).Start()
	return err
}

func StartWith(url string, app string) (err error) {
	err = openWith(url, app).Start()
	return err
}
