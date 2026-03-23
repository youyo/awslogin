package cmd

import (
	"github.com/youyo/awslogin/internal/jsonout"
	"github.com/youyo/awslogin/internal/profile"
)

// ListCmd はプロファイル一覧を表示するコマンド
type ListCmd struct{}

// ListResult は list コマンドの JSON 出力
type ListResult struct {
	Profiles       []profile.Profile      `json:"profiles"`
	CurrentSession *profile.CurrentSession `json:"current_session"`
}

// Run はプロファイル一覧を JSON で出力する
func (c *ListCmd) Run(out *jsonout.Writer) error {
	profiles, err := profile.ListProfiles(profile.DefaultConfigPath())
	if err != nil {
		return err
	}
	if profiles == nil {
		profiles = []profile.Profile{}
	}

	return out.WriteResult(ListResult{
		Profiles:       profiles,
		CurrentSession: profile.DetectCurrentSession(),
	})
}
