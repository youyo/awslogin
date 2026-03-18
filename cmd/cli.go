package cmd

// Globals はすべてのサブコマンドで共有されるフラグとメタ情報を保持する
type Globals struct {
	Open     bool `help:"Open URL in default browser." short:"o"`
	Duration int  `help:"Session duration in seconds." default:"3600" short:"d"`

	// goreleaser ldflags で main.go から注入されるバージョン情報
	AppVersion string `kong:"-"`
	Commit     string `kong:"-"`
	Date       string `kong:"-"`
}

// CLI は Kong のルート構造体
type CLI struct {
	Globals

	Login      LoginCmd      `cmd:"" default:"withargs" help:"Generate AWS console login URL."`
	Version    VersionCmd    `cmd:"" help:"Show version information."`
	Completion CompletionCmd `cmd:"" help:"Generate shell completion script."`
}
