package browse

// Start はデフォルトブラウザで URL を開く
func Start(url string) error {
	return openURL(url).Run()
}
