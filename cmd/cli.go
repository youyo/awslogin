package cmd

// Globals はすべてのサブコマンドで共有されるフラグとメタ情報を保持する
type Globals struct {
	Open     bool `help:"Open URL in default browser (true/false)." short:"o" env:"AWSLOGIN_OPEN"`
	Duration int  `help:"Session duration in seconds (900-43200)." default:"3600" short:"d" env:"AWSLOGIN_DURATION"`

	// goreleaser ldflags で main.go から注入されるバージョン情報
	AppVersion string `kong:"-"`
}

// CLI は Kong のルート構造体
type CLI struct {
	Globals

	Login      LoginCmd      `cmd:"" default:"withargs" help:"Generate AWS console login URL."`
	List       ListCmd       `cmd:"" help:"List available AWS profiles and current session."`
	Version    VersionCmd    `cmd:"" help:"Show version information."`
	Completion CompletionCmd `cmd:"" help:"Generate shell completion script."`
}
